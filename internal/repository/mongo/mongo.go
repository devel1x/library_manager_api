package mongo

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"template/internal/config"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(cfg *config.Mongo) (*mongo.Client, error) {
	fmt.Println(cfg.URI)
	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting mongo: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println(cfg.DBName)
	if err := client.Database(cfg.DBName).RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return client, nil
}
