package router

import (
	"github.com/ilya-korotya/solid/server/handler"
	"github.com/ilya-korotya/solid/server/middleware"
	"github.com/julienschmidt/httprouter"
)

type Config struct {
	Handlers *handler.Handle
}

func (c Config) setBasicRouter(router *httprouter.Router) {
	router.POST("/", middleware.SetCustomContext(c.Handlers.NewUser()))
}

func (c Config) Install() *httprouter.Router {
	router := httprouter.New()
	c.setBasicRouter(router)
	return router
}
