package dto

import (
	v1 "template/internal/delivery/http/v1"
	"template/internal/entity"
	"time"
)

func FormToUser(form v1.UserSignupForm) *entity.User {
	isAdm := 0
	if form.Secret != "" {
		isAdm = 1
	}

	user := entity.User{
		Username:  form.Username,
		CreatedAt: time.Now(),
		Role:      isAdm,
	}
	return &user
}
