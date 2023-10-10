package midlewares

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	auth "github.com/omaily/JWT/internal/jwt"
	libResponse "github.com/omaily/JWT/internal/server/response"
)

func AuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {

		token := request.Header["Authorization"]
		if request.Header["Authorization"] == nil || token[0] == "" {
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain an access token")))
			return
		}
		tokenString := token[0]
		err := auth.ValidateToken(tokenString)
		if err != nil {
			render.Render(write, request, libResponse.ErrInvalidRequest(err))
			return
		}

		next.ServeHTTP(write, request)
	})
}

func AuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {

		token, err := request.Cookie("access_token")
		if err != nil {
			if err == http.ErrNoCookie {
				render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain cookie")))
				return
			}
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain an access token")))
			return
		}

		err = auth.ValidateToken(token.Value)
		if err != nil {
			render.Render(write, request, libResponse.ErrInvalidRequest(err))
			return
		}

		next.ServeHTTP(write, request)
	})
}

func AuthRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {

		token, err := request.Cookie("access_token")
		if err != nil {
			slog.Error(err.Error())
			if err == http.ErrNoCookie {
				render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain cookie")))
				return
			}
			render.Render(write, request, libResponse.ErrInvalidRequest(errors.New("request does not contain an access token")))
			return
		}

		cookieAccesstoken, err := auth.MaintainToken(token.Value)
		if err != nil {
			slog.Error("error maintain token", slog.String("err", err.Error()))
			render.Render(write, request, libResponse.ErrInvalidRequest(err))
			return
		}
		http.SetCookie(write, cookieAccesstoken)

		next.ServeHTTP(write, request)
	})
}
