package server

import (
	"net/http"

	"github.com/ilya-korotya/solid/usecase"
)

// UserCreate register new user in system
func (h *Server) UserCreate(c *Context) error {
	client := &usecase.Client{}
	// TODO: where we have bad err handler
	if err := c.Bind(client); err != nil {
		return err
	}
	/**
	TOOD: but here we have good error handler
	can make api level errors and not usecase?
	*/
	if ok, err := h.Handlers.Register(client); err != nil && !ok {
		return err
	}
	return c.Response(http.StatusOK, map[string]string{"result": "result"})
}

// Users return list of all users available in the system
func (h *Server) Users(c *Context) error {
	users, err := h.Handlers.Users()
	if err != nil {
		return err
	}
	return c.Response(http.StatusOK, users)
}
