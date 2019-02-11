package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ilya-korotya/solid/server/handler"
)

type Server struct {
	Address string
	Port    string
	Handler *handler.Handler
	httprouter.Router
}

func (r *Server) handler(method, pattern string, h handler.Handle) {
	r.Handle(method, pattern, h.ServeHTTP)
}

func (r *Server) Run() {
	r.handler("POST", "/user", r.Handler.CreateUser)
	http.ListenAndServe(r.Address+":"+r.Port, r)
}
