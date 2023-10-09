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

// func (h *Helper) UpdateRefreshToken(rt RT) ([]byte, error) {
// 	defer h.RTCache.Del([]byte(rt.RefreshToken))

// 	userBytes, err := h.RTCache.Get([]byte(rt.RefreshToken))
// 	if err != nil {
// 		return nil, err
// 	}
// 	var u user_service.User
// 	err = json.Unmarshal(userBytes, &u)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return h.GenerateAccessToken(u)
// }

// func (h *Helper) GenerateAccessToken(u user_service.User) ([]byte, error) {

// 	signer, err := jwt.NewSignerHS(jwt.HS256, os.Getenv("JWTKEY"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	builder := jwt.NewBuilder(signer)

// 	claims := UserClaims{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ID:        u.UUID,
// 			Audience:  []string{"users"},
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
// 		},
// 		Email: u.Email,
// 	}
// 	token, err := builder.Build(claims)
// 	if err != nil {
// 		return nil, err
// 	}

// 	h.Logger.Info("create refresh token")
// 	refreshTokenUuid := uuid.New()
// 	userBytes, _ := json.Marshal(u)
// 	err = h.RTCache.Set([]byte(refreshTokenUuid.String()), userBytes, 0)
// 	if err != nil {
// 		h.Logger.Error(err)
// 		return nil, err
// 	}

// 	jsonBytes, err := json.Marshal(map[string]string{
// 		"token":         token.String(),
// 		"refresh_token": refreshTokenUuid.String(),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jsonBytes, nil
// }
