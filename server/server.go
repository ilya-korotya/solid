package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-korotya/solid/entries"

	"github.com/ilya-korotya/solid/usecase"
)

var defaultServer = &Server{
	post: http.NewServeMux(),
	get:  http.NewServeMux(),
}

func InstallDB(db entries.UserStore) {
	defaultServer.DB = db
}

type Handle func(context *Context) error

type Server struct {
	Address  string
	Port     string
	Handlers usecase.UserUsecase
	DB       entries.UserStore
	post     *http.ServeMux
	get      *http.ServeMux
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
	d, _ := json.Marshal(map[string]error{"error": err})
	w.Write(d)
}

func (s *Server) initHandler(h Handle) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(&Context{
			w:  w,
			r:  r,
			DB: s.DB,
		}); err != nil {
			fmt.Println("TODO: procces error:", err)
		}
	}
}

func POST(pattern string, h Handle, server *Server) {
	if server == nil {
		server = defaultServer
	}
	defaultServer.post.HandleFunc(pattern, server.initHandler(h))
}

func GET(pattern string, h Handle, server *Server) {
	if server == nil {
		server = defaultServer
	}
	defaultServer.get.HandleFunc(pattern, server.initHandler(h))
}

func Run(addres string) {
	http.ListenAndServe(addres, defaultServer)
}
