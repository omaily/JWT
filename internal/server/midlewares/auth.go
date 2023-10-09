package midlewares

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	auth "github.com/omaily/JWT/internal/jwt"
	libResponse "github.com/omaily/JWT/internal/server/response"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header["Authorization"]
		if r.Header["Authorization"] == nil || token[0] == "" {
			render.JSON(w, r, libResponse.Error(errors.New("request does not contain an access token")))
			return
		}

		tokenString := token[0]
		err := auth.ValidateToken(tokenString)
		if err != nil {
			render.JSON(w, r, libResponse.Error(err))
			return
		}

		next.ServeHTTP(w, r)
	})
}
