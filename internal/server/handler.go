package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	model "github.com/omaily/JWT/internal/model/user"
	libResponse "github.com/omaily/JWT/internal/server/response"
	"github.com/omaily/JWT/internal/server/validate"
)

func (s *apiServer) helloWorld() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		render.JSON(write, request, "Hello World")
	}
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

func (s *apiServer) authorized(f func(context.Context, *model.User) (string, error)) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {

		user, err := decodeJSON(request)
		if err != nil {
			slog.Error("Failed to decode json", err)
			render.JSON(write, request, libResponse.Error(errors.New("failed to decode request json")))
			return
		}

		valid := validate.ValidateUser(user)
		if valid != nil {
			slog.Error("Failed to validate json", err)
			render.JSON(write, request, valid)
			return
		}

		insertedID, err := f(request.Context(), user)
		if err != nil {
			switch err := err.(type) {
			case *libResponse.InternalError:
				render.JSON(write, request, libResponse.ServerError())
			default:
				render.JSON(write, request, libResponse.Error(err))
			}
			return
		}

		// render.JSON(write, request, api.Ok(insertedID.(primitive.ObjectID).String()))
		render.JSON(write, request, libResponse.Ok(insertedID))

	}
}

func (s *apiServer) createToken() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {

		user, err := decodeJSON(request)
		if err != nil {
			slog.Error("Failed to decode json", err)
			render.JSON(write, request, libResponse.Error(errors.New("failed to decode request json")))
			return
		}

		valid := validate.ValidateUser(user)
		if valid != nil {
			slog.Error("Failed to validate json", err)
			render.JSON(write, request, valid)
			return
		}

		render.JSON(write, request, libResponse.Ok())

	}
}
