package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-korotya/solid/usecase"
	"github.com/lib/pq"
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
	// TODO: add proff and monitor this shit
	switch body.(type) {
	case *pq.Error:
		code = validateDbError(body.(*pq.Error))
	case usecase.CustomError:
		code = validateCustomError(body)
	default:
		code = validateBasicError(body)
	}
	c.w.Header().Set("Content-Type", "application/json")
	d, err := json.Marshal(map[string]string{"error": body.Error()})
	if err != nil {
		return fmt.Errorf("%s:%s", d, err)
	}
	c.w.WriteHeader(code)
	c.w.Write(d)
	return body
}

func validateDbError(err *pq.Error) (code int) {
	switch err.Code {
	case "23505":
		code = http.StatusBadRequest
	}
	return
}

func validateCustomError(err error) (code int) {
	switch usecase.GetType(err) {
	case usecase.BadRequest:
		code = http.StatusBadRequest
	case usecase.NotFound:
		code = http.StatusNotFound
	case usecase.InternalError:
		fallthrough
	default:
		code = http.StatusInternalServerError
	}
	return
}

func validateBasicError(err error) (code int) {
	switch err.(type) {
	case *json.SyntaxError:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	return
}
