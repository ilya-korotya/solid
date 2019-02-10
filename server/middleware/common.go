package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ilya-korotya/solid/usecase"

	"github.com/julienschmidt/httprouter"
)

type CustomContext struct {
	context.Context
	W http.ResponseWriter
}

func SetCustomContext(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := &CustomContext{
			r.Context(),
			w,
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (c *CustomContext) marshal(code int, body interface{}) {
	c.W.WriteHeader(code)
	data, _ := json.Marshal(body)
	// TODO: why didn't this title overwrite?
	c.W.Header().Set("Content-Type", "application/json")
	c.W.Write(data)
}

func (c *CustomContext) HandleError(err error) {
	switch usecase.GetType(err) {
	case usecase.BadRequest:
		c.ResponseError(http.StatusBadRequest, err)
	case usecase.NotFound:
		c.ResponseError(http.StatusNotFound, err)
	default:
		c.ResponseError(http.StatusInternalServerError, err)
	}
}

func (c *CustomContext) ResponseJSON(code int, body interface{}) {
	c.marshal(code, body)
}

func (c *CustomContext) ResponseError(code int, err error) {
	c.marshal(code, map[string]string{"error": err.Error()})
}
