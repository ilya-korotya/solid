package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ilya-korotya/solid/server/context"
	"github.com/ilya-korotya/solid/usecase"
)

type Handle func(c *context.Context) error

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := h(&context.Context{W: w, R: r})
	if err != nil {
		switch usecase.GetType(err) {
		// expected error
		case usecase.InternalError:
			h.proccessError(w, http.StatusInternalServerError, err)
		case usecase.BadRequest:
			h.proccessError(w, http.StatusNotFound, err)
		// fuck! What happened to my service ?!
		case usecase.NotFound:
			h.proccessError(w, http.StatusInternalServerError, err)
		}
	}
}

func (h Handle) proccessError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// TODO: add interface
	w.Write([]byte(err.Error()))
}

type Handler struct {
	user usecase.UserUsecase
}

func New(user usecase.UserUsecase) *Handler {
	return &Handler{
		user: user,
	}
}
