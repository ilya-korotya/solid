package handler

import (
	"net/http"

	"github.com/ilya-korotya/solid/server"
	"github.com/ilya-korotya/solid/usecase"
)

/* TODO: we can implement level errors:
   1. Entries level
   2. Usescase level
   3. Handler level
   4. Database level
   But this is logic stretched across the entire application
**/
func userCreate(c *server.Context) error {
	client := &usecase.Client{}
	if err := c.Bind(client); err != nil {
		return c.ProcessError(err)
	}
	if _, err := c.UserUsecase.Register(client); err != nil {
		return c.ProcessError(err)
	}
	return c.Response(http.StatusOK, client)
}
