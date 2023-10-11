package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	auth "github.com/omaily/JWT/internal/jwt"
	libResponse "github.com/omaily/JWT/internal/model/response"
	"github.com/omaily/JWT/internal/server/midlewares"
)

func RouterSecure(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Use(midlewares.AuthHeader)
		r.Post("/api/work", work())
	})
	router.Get("/api/refresh", refresh())
}

func work() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		render.JSON(write, request, libResponse.Ok("squad went to useful work"))
	}
}

func refresh() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {

		refreshtoken, err := request.Cookie("refresh_token")
		if err != nil {
			slog.Error(err.Error())
			if err == http.ErrNoCookie {
				render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain cookie")))
				return
			}
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain an refresh token")))
			return
		}

		accesstoken := request.Header["Authorization"]
		if request.Header["Authorization"] == nil || accesstoken[0] == "" {
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain an access token")))
			return
		}

		accessrefresh, err := auth.MaintainToken(refreshtoken.Value, accesstoken[0])
		if err != nil {
			slog.Error("error maintain token", slog.String("err", err.Error()))
			render.Render(write, request, libResponse.ErrInvalidRequest(err))
			return
		}

		render.JSON(write, request, libResponse.Bearer(accessrefresh))
	}
}
