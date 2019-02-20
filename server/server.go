package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/ilya-korotya/solid/usecase"
)

var defaultServer = &Server{
	post: http.NewServeMux(),
	get:  http.NewServeMux(),
}

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

func Run(address string, done chan<- struct{}) error {
	server := http.Server{
		Addr:    address,
		Handler: defaultServer,
	}
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("Server off:", err)
		}
		// send default value to 'done' channel
		close(done)
	}()
	return server.ListenAndServe()
}
