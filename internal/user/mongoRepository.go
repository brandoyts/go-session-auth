package user

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const collectionName = "users"

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection(collectionName),
	}
}

func (mr *MongoRepository) FindById(ctx context.Context, id string) (*User, error) {
	// objectId, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, err
	// }

	var found User

	filter := bson.D{{"_id", id}}

	err := mr.collection.FindOne(ctx, filter).Decode(&found)
	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (mr *MongoRepository) FindOne(ctx context.Context, user User) (*User, error) {
	var found User

	err := mr.collection.FindOne(ctx, user).Decode(&found)
	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (mr *MongoRepository) Create(ctx context.Context, user User) (string, error) {
	user.ID = bson.NewObjectID().Hex()

	result, err := mr.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(string), nil
}
