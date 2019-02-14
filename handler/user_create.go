package handler

import (
	"net/http"

	"github.com/ilya-korotya/solid/server"
)

func init() {
	server.POST("/user", userCreate, nil)
	server.GET("/users", users, nil)
}

func userCreate(c *server.Context) error {
	return c.Response(http.StatusCreated, map[string]string{"result": "OK"})
}

func users(c *server.Context) error {
	return c.Response(http.StatusOK, map[string]string{
		"result1": "message1",
		"result2": "message2",
	})
}
