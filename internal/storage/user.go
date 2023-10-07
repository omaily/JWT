package storage

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	model "github.com/omaily/JWT/internal/model/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (st *Storage) identification(ctx context.Context, user *model.User) (*model.User, error) {

	logger := slog.With(
		slog.String("konponent", "storage.user.identification"),
	)

	var result model.User
	err := st.connectUser.FindOne(ctx, bson.M{"_id": user.GUID}).Decode(&result)
	if err != nil {
		logger.Error("query did not match any documents")
		return nil, errors.New("this user not found")
	}

	return &result, nil
}

func (st *Storage) CreateAccount(ctx context.Context, user *model.User) (string, error) {

	logger := slog.With(
		slog.String("konponent", "storage.user.CreateAccount"),
	)

	passwordCript, err := user.SetPassword()
	if err != nil {
		logger.Error("bcrypt library generation error", slog.String("err", err.Error()))
		return "error", errors.New("server internal error")
	}

	result, err := st.connectUser.InsertOne(ctx, &model.User{
		GUID:      primitive.NewObjectID(),
		Email:     user.Email,
		Name:      user.Name,
		Password:  string(passwordCript),
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.Error("error insert", slog.String("err", err.Error()))
		if mongo.IsDuplicateKeyError(err) { // ошибка unique index email_1
			return "error", errors.New("this email is already registered")
		}
		return "error", errors.New("server internal error")
	}

	temp := result.InsertedID.(primitive.ObjectID).Hex()
	logger.Info("inserted user", slog.String("id", temp))
	return temp, nil
}

func (st *Storage) LoginAccount(ctx context.Context, user *model.User) (string, error) {

	logger := slog.With(
		slog.String("konponent", "storage.user.LoginAccount"),
	)

	storedUser, err := st.identification(ctx, user)
	if err != nil {
		return "error", err
	}

	if user.CheckPassword(storedUser.Password) {
		logger.Error("wrong password")
		return "error", errors.New("access denied")
	}

	temp := strconv.FormatBool(true)
	return temp, nil
}
