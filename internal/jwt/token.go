package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type PairToken struct {
	AccessToken  uuid.UUID
	RefreshToken string
}

type Claims struct {
	UUID uuid.UUID `json:"uuid"` // token ID
	GUID string    `json:"id"`   // user ID
	jwt.RegisteredClaims
}

type ModifierToken struct {
	Name    string
	Expires time.Duration
	Key     []byte
}

var (
	// ***Поднять redis или вынести в структуру добавив mutex
	theRedis = make(map[string]PairToken)

	accessKey   = []byte(os.Getenv("ACCESSSECRET"))
	AccessToken = &ModifierToken{
		Name:    "access_token",
		Expires: time.Hour,
		Key:     refreshKey,
	}

	refreshKey   = []byte(os.Getenv("REFRESHSECRET"))
	RefreshToken = &ModifierToken{
		Name:    "refresh_token",
		Expires: 240 * time.Hour,
		Key:     refreshKey,
	}
)
