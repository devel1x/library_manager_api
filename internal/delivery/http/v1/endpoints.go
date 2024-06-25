package v1

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) setRoutes(router chi.Router) {
	router.Route("/v1", func(r chi.Router) {
		r.Route("/user", h.setUserRoutes)
		r.Route("/book", h.setBooksRoutes)
	})
}

func (h *Handler) setUserRoutes(router chi.Router) {
	router.Post("/signup", h.SignUp)
	router.Post("/login", h.Login)

	router.Group(func(r chi.Router) {
		r.Use(h.userIdentity)
		r.Post("/auth/refresh", h.userRefresh)
	})
}

func (h *Handler) setBooksRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(h.adminIdentity)

		r.Post("/", h.CreateBook)
		r.With(h.DeleteBookFromCache).Delete("/{bookISBN}", h.DeleteBookByISBN)
		r.With(h.UpdateBookInCache).Put("/{bookISBN}", h.UpdateBookByISBN)

	})
	router.Group(func(r chi.Router) {
		r.Use(h.userIdentity)

		r.Get("/", h.ListBook)
		r.With(h.FindBookInCache).Get("/{bookISBN}", h.GetBookByISBN)

	})

}
