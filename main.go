package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/omaily/JWT/config"
	model "github.com/omaily/JWT/internal/model/user"
	"github.com/omaily/JWT/internal/server"
	"github.com/omaily/JWT/internal/storage"
)

var (
	logger *slog.Logger
)

func init() {
	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		// 	if a.Key == slog.SourceKey {
		// 		source := a.Value.Any().(*slog.Source)
		// 		source.File = filepath.Base(source.File)
		// 	}
		// 	return a
		// }}),
	)
	slog.SetDefault(logger)
}

func main() {
	conf := config.MustLoad()

	storage, err := storage.NewStorage(&conf.Storage)
	if err != nil {
		logger.Error("could not initialize storage: %w", err)
		return
	}

	insertedID, err := storage.CreateAccount(context.Background(), &model.User{
		Email:        "new@mail.ru",
		Name:         "test",
		Password:     "test",
		Subscription: "random",
	})
	if err != nil {
		logger.Error("error insert")
		return
	}
	logger.Info("insertedID", insertedID)

	serv, err := server.NewServer(&conf.HTTPServer)
	if err != nil {
		logger.Error("could not initialize chi-router: %w", err)
	}

	serv.Start()
}
