package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"template/internal/entity"
	"template/internal/utils"
	"template/pkg/mongodb"
	"time"
)

func (r *MongoRepo) Authenticate(ctx context.Context) (interface{}, error) { return 0, nil }
func (r *MongoRepo) CreateUser(ctx context.Context, user *entity.User) (interface{}, error) {
	res, err := r.usersCollection.InsertOne(ctx, user)
	if err != nil {
		if mongodb.IsDuplicate(err) {
			return 0, utils.ErrUserAlreadyExists
		}
		return 0, err
	}

	return res.InsertedID, nil
}
func (r *MongoRepo) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}
	filter := bson.D{{"_id", objectID}}
	var user entity.User
	err = r.usersCollection.FindOne(ctx, filter).Decode(&user)
	switch {
	case err == nil:
		return &user, nil
	case errors.Is(err, mongo.ErrNoDocuments):
		return nil, utils.ErrNotExist
	default:
		return nil, err
	}
}

func (r *MongoRepo) GetByCredentials(ctx context.Context, username, password string) (entity.User, error) {
	var user entity.User
	if err := r.usersCollection.FindOne(ctx, bson.M{"username": username, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, utils.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}

func (r *MongoRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error) {
	var user entity.User
	if err := r.usersCollection.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, utils.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}

func (r *MongoRepo) SetSession(ctx context.Context, userID primitive.ObjectID, session entity.Session) error {
	_, err := r.usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}
