package midlewares

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	auth "github.com/omaily/JWT/internal/jwt"
	libResponse "github.com/omaily/JWT/internal/server/response"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {

		token := request.Header["Authorization"]
		if request.Header["Authorization"] == nil || token[0] == "" {
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain an access token")))

			return
		}

		tokenString := token[0]
		err := auth.ValidateToken(tokenString)
		if err != nil {
			render.JSON(write, request, libResponse.ErrInvalidRequest(err))
			return
		}

		next.ServeHTTP(write, request)
	})
}
