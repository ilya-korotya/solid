package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ilya-korotya/solid/server/middleware"
	"github.com/ilya-korotya/solid/usecase"
)

func (h *Handle) NewUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: fuck. Why should i always get context ?!
		ctx := r.Context().(*middleware.CustomContext)
		client := &usecase.Client{}
		if err := json.NewDecoder(r.Body).Decode(client); err != nil {
			ctx.HandleError(err)
			return
		}
		if ok, err := h.Usecase.Register(client); err != nil && !ok {
			ctx.HandleError(err)
			return
		}
		ctx.ResponseJSON(http.StatusOK, client)
	})
}
