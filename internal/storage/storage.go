package storage

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/omaily/JWT/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	instance    *mongo.Client
	connectUser *mongo.Collection
}

func NewStorage(cs *config.Storage) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", cs.Role, cs.Pass, cs.Host, cs.Port))
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		slog.Error("unable to create connection pool: %w", err)
		return nil, err
	}

	return &Storage{
		instance:    client,
		connectUser: client.Database("Person").Collection("users"),
	}, nil
}
