package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"template/internal/entity"
	"template/internal/utils"
	"template/pkg/validator"
)

const BookParam = "bookISBN"

type BookInputForm struct {
	ISBN                string   `json:"isbn,omitempty" bson:"_id"`
	Title               string   `json:"title" bson:"title"`
	Publisher           string   `json:"publisher" bson:"publisher"`
	Author              []string `json:"author" bson:"author"`
	validator.Validator `json:"-" bson:"-"`
}

// @Summary Create Book
// @Description Create a new book
// @Tags Book
// @Accept json
// @Produce json
// @Param book body entity.BookForm true "Book form"
// @Success 201 {string} book successfully created "Book created"
// @Failure 400 {object} entity.BookFormError "Invalid input"
// @Failure 500 {string} Internal server error "Internal server error"
// @Security Bearer
// @Router /api/v1/book [post]
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var form BookInputForm

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.responder.WithBadRequest(w, http.StatusText(http.StatusBadRequest))
		return
	}

	id, err := h.booksService.CreateBook(ctx, &form)
	if err != nil {
		if errors.Is(err, utils.InvalidForm) {
			data, err := json.Marshal(form.BookErrors)
			if err != nil {
				h.responder.WithInternalError(w, err.Error())
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(data)
			if err != nil {
				h.responder.WithInternalError(w, "again error userError")
			}
			return
		}
		if errors.Is(err, utils.ErrBookAlreadyExists) {
			h.responder.WithBadRequest(w, err.Error())
			return
		}
		h.responder.WithInternalError(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	h.responder.WithCreated(w, id)
}

// @Summary List Books
// @Description Get a list of books with optional filters
// @Tags Book
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} []entity.Book
// @Failure 400 {string} Invalid query parameters "Invalid query parameters"
// @Failure 500 {string} Internal server error "Internal server error"
// @Security Bearer
// @Router /api/v1/book [get]
func (h *Handler) ListBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	books, err := h.booksService.ListBook(ctx, page, limit)
	if err != nil {
		if errors.Is(err, utils.ErrBadInput) {
			h.responder.WithBadRequest(w, err.Error())
			return
		}
		h.responder.WithInternalError(w, err.Error())
		return
	}

	h.responder.WithOK(w, books)
}

// @Summary Get Book by ISBN
// @Description Get a book by its ISBN
// @Tags Book
// @Accept json
// @Produce json
// @Param bookISBN path string true "Book ISBN"
// @Success 200 {object} entity.Book
// @Failure 404 {string} Book not found "Book not found"
// @Failure 500 {string} Internal server error "Internal server error"
// @Security Bearer
// @Router /api/v1/book/{bookISBN} [get]
func (h *Handler) GetBookByISBN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, BookParam)
	fmt.Println(2)
	if idParam == "" {
		h.responder.WithBadRequest(w, "empty id param")
		return
	}
	book, err := h.booksService.GetBookByISBN(ctx, idParam)
	if err != nil {
		if errors.Is(err, utils.ErrNotExist) {
			h.responder.WithNotFound(w, "book not found")
			return
		}
		h.responder.WithInternalError(w, err.Error())
		return
	}

	v, ok := ctx.Value(BookParam).(*entity.Book)
	if !ok {
		h.responder.WithInternalError(w, "cannot convert to book")
		return
	}
	*v = *book
	_ = v
	h.responder.WithOK(w, book)
}

// @Summary Update Book by ISBN
// @Description Update an existing book by its ISBN
// @Tags Book
// @Accept json
// @Produce json
// @Param bookISBN path string true "Book ISBN"
// @Param book body entity.BookForm true "Book form"
// @Success 200 {string} book successfully updated "Book updated"
// @Failure 400 {object} entity.BookFormError "Invalid input"
// @Failure 404 {string} book not found "Book not found"
// @Failure 500 {string} Internal server error "Internal server error"
// @Security Bearer
// @Router /api/v1/book/{bookISBN} [put]
func (h *Handler) UpdateBookByISBN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, BookParam)
	if idParam == "" {
		h.responder.WithBadRequest(w, "empty id param")
		return
	}

	var form BookInputForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.responder.WithBadRequest(w, http.StatusText(http.StatusBadRequest))
		return
	}
	form.ISBN = idParam
	book, err := h.booksService.UpdateBookByISBN(ctx, &form)
	if err != nil {
		if errors.Is(err, utils.ErrNotExist) {
			h.responder.WithNotFound(w, "book not found")
			return
		}

	}
	v, ok := ctx.Value(BookParam).(*entity.BookForm)
	if !ok {
		h.responder.WithInternalError(w, "cannot convert to book")
		return
	}
	*v = *book
	_ = v
	h.responder.WithOK(w, "book updated successfully")
}

// @Summary Delete Book by ISBN
// @Description Delete a book by its ISBN
// @Tags Book
// @Accept json
// @Produce json
// @Param bookISBN path string true "Book ISBN"
// @Success 204 {string} book successfully deleted "Book deleted"
// @Failure 404 {string} book not found "Book not found"
// @Failure 500 {string} Internal server error "Internal server error"
// @Security Bearer
// @Router /api/v1/book/{bookISBN} [delete]
func (h *Handler) DeleteBookByISBN(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, BookParam)
	if idParam == "" {
		h.responder.WithBadRequest(w, "empty id param")
		return
	}
	err := h.booksService.DeleteBookByISBN(ctx, idParam)
	if err != nil {
		if errors.Is(err, utils.ErrNotExist) {
			h.responder.WithNotFound(w, "book not found")
			return
		}
		h.responder.WithInternalError(w, err.Error())
		return
	}
	ctx = context.WithValue(ctx, BookParam, idParam)
	h.responder.WithOK(w, "book successfully deleted")
}
