package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDb(database string, url string) (*mongo.Database, error) {
	credentials := options.Credential{
		AuthSource: "admin",
		Username:   "rootusername",
		Password:   "rootpassword",
	}

	clientOptions := options.Client().ApplyURI(url).SetAuth(credentials)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(database), nil
}
