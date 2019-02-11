package handler

import (
	"net/http"

	"github.com/ilya-korotya/solid/server/context"
	"github.com/ilya-korotya/solid/usecase"
)

// CreateUser call usecase rgister via http
func (h *Handler) CreateUser(c *context.Context) error {
	client := &usecase.Client{}
	if err := c.RequestBody(client); err != nil {
		return err
	}
	if ok, err := h.user.Register(client); !ok && err != nil {
		return err
	}
	return c.Response(http.StatusOK, client)
}
