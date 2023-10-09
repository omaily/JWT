package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/omaily/JWT/internal/server/midlewares"
	libResponse "github.com/omaily/JWT/internal/server/response"
)

func RouterSecure(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Use(midlewares.BasicAuth)
		r.Post("/api/order", order())
	})
}

func order() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		render.JSON(write, request, libResponse.Ok("anti-corruption squad sent"))
	}
}
