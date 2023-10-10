package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey               = []byte(os.Getenv("JWTKEY"))
	AccessTokenLifetime  = 2 * time.Hour
	RefreshTokenLifetime = 72 * time.Hour
)

type Claims struct {
	GUID     string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateToken(id string, email string, username string) (*http.Cookie, error) {
	expirationTime := time.Now().Add(AccessTokenLifetime)
	jwtExpirationTime := jwt.NewNumericDate(expirationTime)
	claims := &Claims{
		GUID:     id,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwtExpirationTime,
			Issuer:    "my crew, my gang, my family",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("сгенерирован новый токен")
	return &http.Cookie{
		Name:    "access_token",
		Path:    "/",
		Value:   tokenString,
		Expires: expirationTime,
	}, nil
}

func ValidateToken(tokenArrived string) error {
	logger := slog.With(
		slog.String("konponent", "jwt.ValidateToken"),
	)
	if _, err := checkSignature(tokenArrived); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func MaintainToken(tokenArrived string) (*http.Cookie, error) {

	logger := slog.With(
		slog.String("konponent", "jwt.MaintainToken"),
	)

	claims, err := checkSignature(tokenArrived)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	realLifeTimeToken := AccessTokenLifetime - time.Until(claims.ExpiresAt.Time)
	if realLifeTimeToken < time.Minute {
		err := errors.New("too little time has passed since the token was created")
		logger.Error(err.Error())
		return nil, err
	}
	return GenerateToken(claims.GUID, claims.Email, claims.Username)
}

func checkSignature(tokenArrived string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenArrived, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token expired")
	}
}
