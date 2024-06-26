package bookService

import (
	"context"
	"strconv"
	v1 "template/internal/delivery/http/v1"
	"template/internal/dto"
	"template/internal/entity"
	"template/internal/utils"
	"template/pkg/validator"
	"unicode/utf8"
)

const (
	limitDefault = 5
	pageDefault  = 1
)

type BookService struct {
	bookRepo bookRepo
}

func NewBookService(bookRepo bookRepo) *BookService {
	return &BookService{
		bookRepo: bookRepo,
	}
}

type bookRepo interface {
	GetBookByISBN(ctx context.Context, id string) (*entity.Book, error)
	CreateBook(ctx context.Context, book *entity.BookFormCreate) (interface{}, error)
	ListBook(ctx context.Context, page, pageSize int) (*entity.PaginatedBooks, error)
	UpdateBookByISBN(ctx context.Context, book *entity.BookFormUpdate) error
	DeleteBookByISBN(ctx context.Context, id string) error
}

func (s *BookService) GetBookByISBN(ctx context.Context, id string) (*entity.Book, error) {

	return s.bookRepo.GetBookByISBN(ctx, id)
}

func (s *BookService) ListBook(ctx context.Context, pageStr, limitStr string) (*entity.PaginatedBooks, error) {
	var page, limit int
	var err error
	if utf8.RuneCountInString(pageStr) == 0 {
		page = pageDefault
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, utils.ErrBadInput
		}
		if page <= 0 {
			return nil, utils.ErrBadInput
		}
	}

	if utf8.RuneCountInString(limitStr) == 0 {
		limit = limitDefault
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, utils.ErrBadInput
		}
		if limit <= 0 {
			return nil, utils.ErrBadInput
		}
	}

	return s.bookRepo.ListBook(ctx, page, limit)
}

func (s *BookService) UpdateBookByISBN(ctx context.Context, book *v1.BookInputForm) (*entity.BookFormUpdate, error) {
	form := dto.BookToFormUpdate(book)
	return form, s.bookRepo.UpdateBookByISBN(ctx, form)
}

func (s *BookService) DeleteBookByISBN(ctx context.Context, id string) error {

	return s.bookRepo.DeleteBookByISBN(ctx, id)
}

func (s *BookService) CreateBook(ctx context.Context, book *v1.BookInputForm) (interface{}, error) {
	book.CheckField(validator.CheckISBN(book.ISBN, 13), &book.BookErrors.ISBN, "ISBN must be exactly 13 chars long")
	book.CheckField(validator.NotBlank(book.Title), &book.BookErrors.Title, "This field cannot be blank")
	book.CheckField(validator.NotBlank(book.Publisher), &book.BookErrors.Publisher, "This field cannot be blank")
	book.CheckField(validator.CheckArr(book.Author), &book.BookErrors.Author, "This field cannot be blank")

	if !book.ValidBook() {
		return nil, utils.InvalidForm
	}

	form := dto.BookToForm(book)

	return s.bookRepo.CreateBook(ctx, form)
}
