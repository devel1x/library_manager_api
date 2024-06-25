package mongo

import "go.mongodb.org/mongo-driver/mongo"

type MongoRepo struct {
	usersCollection *mongo.Collection
	booksCollection *mongo.Collection
}

func NewRepoMongo(usersCollection *mongo.Collection, booksCollection *mongo.Collection) *MongoRepo {
	return &MongoRepo{
		usersCollection: usersCollection,
		booksCollection: booksCollection,
	}
}
