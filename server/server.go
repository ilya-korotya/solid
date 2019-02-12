package server

import (
	"encoding/json"
	"net/http"

	"github.com/ilya-korotya/solid/usecase"
)

func proccesError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	d, _ := json.Marshal(map[string]error{"error": err})
	w.Write(d)
}

type Handle func(context *Context) error

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(&Context{w: w, r: r})
	if err != nil {
		// TODO: implement full handle error
		switch usecase.GetType(err) {
		case usecase.BadRequest:
			proccesError(w, http.StatusBadRequest, err)
		}
	}
}

type Server struct {
	Address  string
	Port     string
	Handlers usecase.UserUsecase
	post     http.ServeMux
	get      http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.get.ServeHTTP(w, r)
	case "POST":
		s.post.ServeHTTP(w, r)
	}
}

func (r *Server) POST(pattern string, h Handle) {
	r.post.Handle(pattern, h)
}

func (r *Server) GET(pattern string, h Handle) {
	r.get.Handle(pattern, h)
}

func (r *Server) Run() {
	r.POST("/user", r.UserCreate)
	r.GET("/users", r.Users)
	http.ListenAndServe(r.Address+":"+r.Port, r)
}
