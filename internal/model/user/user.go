package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	GUID         primitive.ObjectID `json:"id" bson:"_id"`
	Email        string             `json:"email" bson:"email" validate:"required"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	Password     string             `json:"password" bson:"password" validate:"required"`
	Subscription string             `bson:"subscription"`
	CreatedAt    time.Time          `bson:"created_at"`
}

func (u *User) SetPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(passwordCript string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordCript), []byte(u.Password)) != nil
}
