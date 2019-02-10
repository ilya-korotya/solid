package server

import (
	"net/http"
)

type Config struct {
	Address string
	Port    string
}

func (c *Config) Run(router http.Handler) error {
	return http.ListenAndServe(c.Address+":"+c.Port, router)
}
