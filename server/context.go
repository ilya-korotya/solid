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

// NewContext create context with custom writer and reader
// TODO: may be better write w and r as public and remove the constructor
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w: w,
		r: r,
	}
}

// Bind parse body to target structure
func (c *Context) Bind(target interface{}) error {
	defer c.r.Body.Close()
	return json.NewDecoder(c.r.Body).Decode(target)
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

// ProcessError proccess error and proxy it to loger
func (c *Context) ProcessError(body error) error {
	var code int
	c.w.Header().Set("Content-Type", "application/json")
	switch usecase.GetType(body) {
	case usecase.BadRequest:
		code = http.StatusBadRequest
	case usecase.NotFound:
		code = http.StatusNotFound
	case usecase.InternalError:
		fallthrough
	default:
		code = http.StatusInternalServerError
	}
	d, err := json.Marshal(map[string]string{"error": body.Error()})
	if err != nil {
		return err
	}
	c.w.WriteHeader(code)
	c.w.Write(d)
	return body
}
