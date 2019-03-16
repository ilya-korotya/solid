package handler

import (
	"net/http"
	"time"

	"github.com/ilya-korotya/solid/server"
)

func users(c *server.Context) error {
	users, err := c.UserUsecase.Users()
	if err != nil {
		return c.ProcessError(err)
	}
	return c.Response(http.StatusOK, users)
}

// use for check graceful shutdown
func foo(c *server.Context) error {
	time.Sleep(30 * time.Second)
	return c.Response(http.StatusOK, map[string]string{"message": "OK"})
}
