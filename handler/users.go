package handler

import (
	"net/http"

	"github.com/ilya-korotya/solid/server"
)

func users(c *server.Context) error {
	users, err := c.UserUsecase.Users()
	if err != nil {
		return err
	}
	return c.Response(http.StatusOK, users)
}
