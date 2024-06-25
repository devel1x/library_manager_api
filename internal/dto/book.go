package dto

import (
	v1 "template/internal/delivery/http/v1"
	"template/internal/entity"
	"time"
)

func BookToForm(book *v1.BookInputForm) *entity.BookForm {
	form := entity.BookForm{
		ISBN:      book.ISBN,
		Title:     book.Title,
		Publisher: book.Publisher,
		Author:    book.Author,
		CreatedAt: time.Now(),
	}

	return &form
}

func BookToFormUpdate(book *v1.BookInputForm) *entity.BookForm {
	form := entity.BookForm{
		ISBN:      book.ISBN,
		Title:     book.Title,
		Publisher: book.Publisher,
		Author:    book.Author,
		UpdatedAt: time.Now(),
	}

	return &form
}
