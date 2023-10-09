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
)

func Router(router *chi.Mux, storage *storage.Storage) {
	router.Route("/api/auth", func(r chi.Router) {
		r.Post("/createAccount", authorized(storage.CreateAccount))
		r.Post("/login", authorized(storage.LoginAccount))
		r.Post("/createToken", createToken(storage.LoginAccount))
	})
}

func decodeJSON(request *http.Request) (*model.User, error) {
	var user model.User
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields() // генерирует ошибку если в json есть поля которых нет в структуре
	err := decoder.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func authorized(f func(context.Context, *model.User) (string, error)) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {

		user, err := decodeJSON(request)
		if err != nil {
			slog.Error("Failed to decode json", slog.String("err", err.Error()))
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("failed to decode request json")))
			return
		}

		valid := validate.ValidateUser(user)
		if valid != nil {
			slog.Error("Failed to validate json")
			render.Render(write, request, valid)
			return
		}

		insertedID, err := f(request.Context(), user)
		if err != nil {
			render.Render(write, request, libResponse.ErrReview(err))
			return
		}

		render.JSON(write, request, libResponse.Ok(insertedID))
	}
}

func createToken(f func(context.Context, *model.User) (string, error)) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {

		user, err := decodeJSON(request)
		if err != nil {
			slog.Error("Failed to decode json", err)
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("failed to decode request json")))
			return
		}

		valid := validate.ValidateUser(user)
		if valid != nil {
			slog.Error("Failed to validate json")
			render.Render(write, request, valid)
			return
		}

		token, err := f(request.Context(), user)
		if err != nil {
			render.Render(write, request, libResponse.ErrInvalidRequest(err))
			return
		}

		render.JSON(write, request, libResponse.Bearer(token))
	}
}
