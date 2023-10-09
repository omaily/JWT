package auth

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	model "github.com/omaily/JWT/internal/model/user"
)

var jwtKey = []byte(os.Getenv("JWTKEY"))

type Claims struct {
	Email    string `json:"email"`
	Username string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateToken(u *model.User) (string, error) {

	expirationTime := jwt.NewNumericDate(time.Now().Add(time.Hour * 72))
	claims := &Claims{
		Email:    u.Email,
		Username: u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
			Issuer:    "my crew, my gang, my family",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) error {
	logger := slog.With(
		slog.String("konponent", "jwt.ValidateToken"),
	)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		return nil
	} else {
		logger.Error(err.Error())
		return errors.New("token expired")
	}
}
