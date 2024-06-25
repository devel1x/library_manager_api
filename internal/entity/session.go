package entity

import "time"

type Session struct {
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" bson:"expires_at"`
}

type Tokens struct {
	AccessToken  string `json:"access_token" bson:"refresh_token"`
	RefreshToken string `json:"refresh_token" bson:"expires_at"`
}

type RefreshInput struct {
	Token string `json:"token" binding:"required"`
}
