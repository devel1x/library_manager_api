package v1

import (
	"context"
	"template/internal/entity"
)

type userService interface {
	Login(ctx context.Context, input *UserLoginForm) (entity.Tokens, error)
	SignUp(ctx context.Context, form *UserSignupForm) (interface{}, error)
	GetUserByID(id string) (*entity.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (entity.Tokens, error)
}

type bookService interface {
	ListBook(ctx context.Context, page, limit string) (*entity.PaginatedBooks, error)
	GetBookByISBN(ctx context.Context, id string) (*entity.Book, error)
	UpdateBookByISBN(ctx context.Context, book *BookInputForm) (*entity.BookFormUpdate, error)
	DeleteBookByISBN(ctx context.Context, id string) error
	CreateBook(ctx context.Context, book *BookInputForm) (interface{}, error)
}

type redisInterface interface {
	InsertBook(ctx context.Context, book *entity.Book) error
	FindBookByISBN(ctx context.Context, id string) (*entity.Book, error)
	DeleteBookByISBN(ctx context.Context, id string) error
	UpdateBookByISBN(ctx context.Context, book *entity.BookFormUpdate) error
}
