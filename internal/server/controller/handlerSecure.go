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
		// r.Use(midlewares.AuthHeader)
		r.Use(midlewares.AuthCookie)
		r.Post("/api/order", order())
	})

	router.Group(func(r chi.Router) {
		r.Use(midlewares.AuthRefresh)
		r.Get("/api/maintain", refresh())
	})
}

func order() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		render.JSON(write, request, libResponse.Ok("squad left order"))
	}
}

func refresh() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		render.JSON(write, request, libResponse.Ok("squad left refresh"))
	}
}
