package validate

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	model "github.com/omaily/JWT/internal/model/user"
	libResponse "github.com/omaily/JWT/internal/server/response"
)

func ValidateUser(user *model.User) *libResponse.ErrResponse {

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// Нужен только первый тег json (игнорируем omitempty)
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(user)
	if err != nil {
		slog.Error("invalid request", err)
		return ValidateError(err.(validator.ValidationErrors))
	}
	return nil
}

func ValidateError(e validator.ValidationErrors) *libResponse.ErrResponse {

	res := libResponse.ErrValidaete(errors.New("failed validete structure"))

	for _, err := range e { // ошибка может прийти не по одному полю
		e := libResponse.ValidateError{
			NameStruct:    err.StructField(),
			Type:          fmt.Sprintf("%v", err.Type()),
			NameFieldJson: err.Field(),
			Value:         fmt.Sprintf("%v", err.Value()),
			Message:       err.Error(),
		}
		res.Valid = append(res.Valid, e)
	}

	return res
}
