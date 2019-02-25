package server

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilya-korotya/solid/usecase"
)

var (
	fileDescriptorFlag = flag.Int("fd", 0, "File discriptor for http listener")
	defaultServer      = &Server{
		post: http.NewServeMux(),
		get:  http.NewServeMux(),
	}
)

func InstallUserUsecase(uc usecase.UserUsecase) {
	defaultServer.userUsecase = uc
}

type Handle func(context *Context) error

type Server struct {
	userUsecase usecase.UserUsecase
	post        *http.ServeMux
	get         *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.get.ServeHTTP(w, r)
	case "POST":
		s.post.ServeHTTP(w, r)
	}
}

func proccesError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	d, _ := json.Marshal(map[string]string{"error": err.Error()})
	w.Write(d)
}

func (s *Server) initHandler(h Handle) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(&Context{
			w:           w,
			r:           r,
			UserUsecase: s.userUsecase,
		}); err != nil {
			switch usecase.GetType(err) {
			case usecase.BadRequest:
				proccesError(w, http.StatusBadRequest, err)
			case usecase.NotFound:
				proccesError(w, http.StatusNotFound, err)
			case usecase.InternalError:
				fallthrough
			default:
				proccesError(w, http.StatusInternalServerError, err)
			}
		}
	}
}

func POST(pattern string, h Handle) {
	defaultServer.post.HandleFunc(pattern, defaultServer.initHandler(h))
}

func GET(pattern string, h Handle) {
	defaultServer.get.HandleFunc(pattern, defaultServer.initHandler(h))
}

// Run http server with graceful restart and shutdow
// Use signal SIGINT for graceful shutdow
// Use signal SIGTSTP for graceful restart
func Run(address string, done chan<- struct{}) error {
	var parentListener net.Listener
	flag.Parse()
	server := http.Server{
		Addr:    address,
		Handler: defaultServer,
	}
	go func() {
		sigtstp := make(chan os.Signal, 1)
		signal.Notify(sigtstp, syscall.SIGTSTP)
		<-sigtstp
		childListener, ok := parentListener.(*net.TCPListener)
		if !ok {
			panic("cannot cast file descriptor to TCPListener")
		}
		// get file by listener
		file2, err := childListener.File()
		if err != nil {
			panic("cannot get file via descriptor")
		}
		// get file descriptor
		fd1 := int(file2.Fd())
		// make copy this descriptor without FD_CLOEXEC flag
		fd2, err := syscall.Dup(fd1)
		if err != nil {
			panic("cannot create of copy without FD_CLOEXEC flag")
		}
		// run forke proccess with custom file descriptor
		cmd := exec.Command("./solid", fmt.Sprint("--fd=", fd2))
		if err := cmd.Start(); err != nil {
			panic(fmt.Sprintln("cannot run forke procces:", err))
		}
		server.Shutdown(context.Background())
		// wait for a copy of the process to start
		time.Sleep(10 * time.Second)
		parentListener.Close()
		close(done)
	}()
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("Server off:", err)
		}
		close(done)
	}()
	var err error
	// use the descriptor number to create a listener from this descriptor
	if *fileDescriptorFlag != 0 {
		// create file from file descriptor
		fileListen := os.NewFile(uintptr(*fileDescriptorFlag), "parrent")
		parentListener, err = net.FileListener(fileListen)
		if err != nil {
			return err
		}
	} else {
		// create default tcp listener
		parentListener, err = net.Listen("tcp", address)
		if err != nil {
			return err
		}
	}
	return server.Serve(parentListener)
}
