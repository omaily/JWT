package storage

import (
	"context"
	"errors"
	"log/slog"
	"time"

	model "github.com/omaily/JWT/internal/model/user"
	libResponse "github.com/omaily/JWT/internal/server/response"

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
		return "err", &libResponse.InternalError{}
	}

	result, err := st.connectUser.InsertOne(ctx, &model.User{
		GUID:      primitive.NewObjectID(),
		Email:     user.Email,
		Name:      user.Name,
		Password:  string(passwordCript),
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.Error("mongoDB error insert", slog.String("err", err.Error()))
		if mongo.IsDuplicateKeyError(err) { // ошибка unique index email_1
			return "err", errors.New("this email is already registered")
		}
		return "err", &libResponse.InternalError{}
	}

	temp := result.InsertedID.(primitive.ObjectID).Hex()
	logger.Info("inserted user", slog.String("id", temp))
	return temp, nil
}

func (st *Storage) LoginAccount(ctx context.Context, user *model.User) error {

	logger := slog.With(
		slog.String("konponent", "storage.user.LoginAccount"),
	)

	storedUser, err := st.identification(ctx, user)
	if err != nil {
		return err
	}

	if user.CheckPassword(storedUser.Password) {
		logger.Error("wrong password")
		return errors.New("access denied")
	}

	// tokenString, err := auth.GenerateToken(user)
	// if err != nil {
	// 	logger.Error("error creating token", slog.String("err", err.Error()))
	// 	return "err", &libResponse.InternalError{}
	// }

	return nil
}
