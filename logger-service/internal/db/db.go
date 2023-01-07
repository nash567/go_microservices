package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConnection(config *Config) (*mongo.Client, error) {
	fmt.Println("username is", config.UserName, config.Password)
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	clientOptions.SetAuth(options.Credential{
		Username: config.UserName,
		Password: config.Password,
	})
	conn, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	return conn, nil
}

// Verify ensures the connection is available for use.
func Verify(conn *mongo.Client) error {

	if err := conn.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
