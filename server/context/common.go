package context

import (
	"encoding/json"
	"net/http"
)

// Context request with additional utilities
type Context struct {
	W http.ResponseWriter
	R *http.Request
}

// RequestBody parse body to structure
// TODO: super duper mega shit
func (c *Context) RequestBody(targger interface{}) error {
	defer c.R.Body.Close()
	return json.NewDecoder(c.R.Body).Decode(targger)
}

// Response sends data to client in JSON type
func (c *Context) Response(code int, body interface{}) error {
	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(code)
	d, err := json.Marshal(body)
	if err != nil {
		return err
	}
	c.W.Write(d)
	return nil
}
