package server

import (
	"encoding/json"
	"net/http"

	"github.com/ilya-korotya/solid/usecase"
)

// Context request with additional utilities
type Context struct {
	w           http.ResponseWriter
	r           *http.Request
	UserUsecase usecase.UserUsecase
}

// Bind parse body to targger structure
func (c *Context) Bind(targger interface{}) error {
	defer c.r.Body.Close()
	return json.NewDecoder(c.r.Body).Decode(targger)
}

// Response sends data to client in JSON type
func (c *Context) Response(code int, body interface{}) error {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(code)
	d, err := json.Marshal(body)
	if err != nil {
		return err
	}
	c.w.Write(d)
	return nil
}
