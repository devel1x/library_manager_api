package entity

import (
	"time"
)

type Book struct {
	ISBN      string    `json:"isbn" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Publisher string    `json:"publisher" bson:"publisher"`
	Author    []string  `json:"author" bson:"author"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type BookFormCreate struct {
	ISBN      string    `json:"isbn" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Publisher string    `json:"publisher" bson:"publisher"`
	Author    []string  `json:"author" bson:"author"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type BookFormUpdate struct {
	ISBN      string    `json:"isbn" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Publisher string    `json:"publisher" bson:"publisher"`
	Author    []string  `json:"author" bson:"author"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type BookFormError struct {
	ISBN      string `json:"isbn,omitempty" bson:"isbn,omitempty"`
	Title     string `json:"title,omitempty" bson:"title,omitempty"`
	Publisher string `json:"publisher,omitempty" bson:"publisher,omitempty"`
	Author    string `json:"author,omitempty" bson:"author,omitempty"`
}

type PaginatedBooks struct {
	Books    []*Book `json:"books"`
	LastPage int     `json:"last_page"`
}
