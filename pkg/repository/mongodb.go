package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	Host       string
	Port       int
	Database   string
	Collection string
}

func NewMongoDB(cfg *Config) (*mongo.Collection, error) {
	dbOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port))
	client, err := mongo.NewClient(dbOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.Database).Collection(cfg.Collection), nil
}
