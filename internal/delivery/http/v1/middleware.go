package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"template/internal/entity"
	"template/internal/utils"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) FindBookInCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		bookISBN := chi.URLParam(r, BookParam)
		bookISBN = strings.TrimSpace(bookISBN)
		book, err := h.cache.FindBookByISBN(ctx, bookISBN)

		if err != nil {
			if errors.Is(err, utils.ErrNotExist) {
				r = r.WithContext(context.WithValue(r.Context(), BookParam, &entity.Book{}))
				next.ServeHTTP(w, r)
				v, ok := r.Context().Value(BookParam).(*entity.Book)
				if !ok {
					h.responder.WithInternalError(w, "couldn't convert book to struct")
					return
				}
				err = h.cache.InsertBook(ctx, v)
				if err != nil {
					h.responder.WithInternalError(w, "error inserting book in cache")
				}
				return
			}
			h.responder.WithInternalError(w, fmt.Sprintf("error getting book from cache: %v", err))
			return
		}
		h.responder.WithOK(w, book)
	})
}

func (h *Handler) UpdateBookInCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), BookParam, &entity.BookFormCreate{}))
		next.ServeHTTP(w, r)
		v, ok := r.Context().Value(BookParam).(*entity.BookFormUpdate)
		if !ok {
			h.responder.WithInternalError(w, "couldnt convert book to struct")
			return
		}

		ctx := context.Background()
		err := h.cache.UpdateBookByISBN(ctx, v)
		if err != nil {
			if errors.Is(err, utils.ErrNotExist) {
				return
			}
			h.responder.WithInternalError(w, fmt.Sprintf("error updating book in cache: %v", err))
		}
	})
}

func (h *Handler) DeleteBookFromCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		ctx := r.Context()
		isbn := chi.URLParam(r, BookParam)
		if isbn == "" {
			return
		}
		err := h.cache.DeleteBookByISBN(ctx, isbn)
		if err != nil {
			if errors.Is(err, utils.ErrNotExist) {
				return
			}
			h.responder.WithInternalError(w, fmt.Sprintf("error getting book from cache: %v", err))
		}
	})
}

func (h *Handler) adminIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := h.parseAuthHeader(r)
		if err != nil {
			h.responder.WithUnauthorizedError(w)
			return
		}
		fmt.Println(id)
		user, err := h.userService.GetUserByID(id)
		if err != nil {
			if errors.Is(err, utils.ErrUserNotFound) {
				h.responder.WithUnauthorizedError(w)
				return
			}
			h.responder.WithInternalError(w, "error authorizing user")
			return
		}
		// ----------
		if user.Role != 1 {
			h.responder.WithForbiddenError(w)
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := h.parseAuthHeader(r)
		if err != nil {
			h.responder.With(http.StatusUnauthorized, w, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) parseAuthHeader(r *http.Request) (string, error) {
	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
	// return true false, check admin in db
}
