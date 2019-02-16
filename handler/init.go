package handler

import "github.com/ilya-korotya/solid/server"

// This package is used to install all paths
func init() {
	server.POST("/user", userCreate)
	server.GET("/users", users)
}
