package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	HashedPassword string             `json:"password" bson:"password"`
	Role           int                `json:"role" bson:"role"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
}

type UserFormError struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Secret   string `json:"secret,omitempty" bson:"secret,omitempty"`
}
