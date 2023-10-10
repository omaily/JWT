package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	model "github.com/omaily/JWT/internal/model/user"
	libResponse "github.com/omaily/JWT/internal/server/response"
	"github.com/omaily/JWT/internal/server/validate"
	"github.com/omaily/JWT/internal/storage"

	auth "github.com/omaily/JWT/internal/jwt"
)

type userI interface {
	CreateAccount(context.Context, *model.User) (string, error)
	LoginAccount(context.Context, *model.User) error
}

func Router(router *chi.Mux, storage *storage.Storage) {
	router.Route("/api", func(r chi.Router) {
		r.Post("/createAccount", authorized(storage))
		r.Post("/login", login(storage))
	})
}

func checkRequestJson(write http.ResponseWriter, request *http.Request) (*model.User, error) {
	var user model.User
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields() //генерирует ошибку если в json есть поля которых нет в структуре
	err := decoder.Decode(&user)
	if err != nil {
		slog.Error("Failed to decode json", slog.String("err", err.Error()))
		render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("failed to decode json")))
		return nil, err
	}

	valid := validate.ValidateUser(&user)
	if valid != nil {
		slog.Error("Failed to validate json")
		render.Render(write, request, valid)
		return nil, err
	}

	return &user, nil
}

func authorized(u userI) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {

		user, err := checkRequestJson(write, request)
		if err != nil {
			return
		}

		insertedID, err := u.CreateAccount(request.Context(), user)
		if err != nil {
			render.Render(write, request, libResponse.ErrReview(err))
			return
		}

		render.JSON(write, request, libResponse.Ok(insertedID))
	}
}

func login(u userI) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {

		user, err := checkRequestJson(write, request)
		if err != nil {
			return
		}

		if err := u.LoginAccount(request.Context(), user); err != nil {
			render.Render(write, request, libResponse.ErrInvalidRequest(err))
			return
		}

		cookieAccesstoken, err := auth.GenerateToken(user.GUID.Hex(), user.Email, user.Name)
		if err != nil {
			slog.Error("error creating token", slog.String("err", err.Error()))
			render.Render(write, request, libResponse.ErrInvalidRequest(err))
			return
		}
		http.SetCookie(write, cookieAccesstoken)

		// http.SetCookie(write, &http.Cookie{
		// 	Name:  "refresh_token",
		// 	Path:  "/",
		// 	Value: "tak.tuk.tik",
		// })

		render.JSON(write, request, libResponse.Ok())
	}
}
