package v1

import (
	httpSwagger "github.com/swaggo/http-swagger"
	_ "template/docs"
	"template/internal/delivery/http"
	"template/pkg/auth"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	responder    http.Responder
	logger       *zap.Logger
	userService  userService
	booksService bookService
	cache        redisInterface
	tokenManager auth.TokenManager
}

func NewHandler(
	responder http.Responder,
	logger *zap.Logger,
	userService userService,
	booksService bookService,
	cache redisInterface,
	manager auth.TokenManager,
) *Handler {
	return &Handler{
		responder:    responder,
		logger:       logger,
		userService:  userService,
		booksService: booksService,
		cache:        cache,
		tokenManager: manager,
	}
}

func SetHandler(
	mux *chi.Mux,
	responder http.Responder,
	logger *zap.Logger,
	userService userService,
	booksService bookService,
	cache redisInterface,
	manager auth.TokenManager,
) {
	handler := NewHandler(responder, logger, userService, booksService, cache, manager)
	mux.Route("/api", handler.setRoutes)
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
