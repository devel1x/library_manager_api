package server

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	_ "template/docs"
	"template/internal/config"
	http2 "template/internal/delivery/http"
	"template/internal/delivery/http/v1"
	db "template/internal/repository/mongo"
	cache "template/internal/repository/redis"
	book_service "template/internal/service/book"
	user_service "template/internal/service/user"
	"template/pkg/auth"
	"template/pkg/hash"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	cfg    *config.Config
	router *chi.Mux
	logger *zap.Logger
	db     *mongo.Client //iocloser
	cache  *redis.Client
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Initialize() error {
	var (
		err error
	)

	if err != nil {
		return err
	}
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339TimeEncoder

	output := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			CallerKey:      "caller",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	core := zapcore.NewTee(output)
	a.logger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		},
		)),
	)
	a.router = chi.NewRouter()
	a.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: a.cfg.App.Cors.AllowOrigins,
		AllowedMethods: a.cfg.App.Cors.AllowMethods,
		AllowedHeaders: a.cfg.App.Cors.AllowHeaders,
	}))
	if err = a.setHandler(); err != nil {
		a.logger.Info(err.Error())
		return err
	}
	return nil
}

func (a *App) setHandler() error {
	var err error

	a.db, err = db.NewMongoDB(&a.cfg.Repository.Mongo)
	if err != nil {
		return err
	}

	a.cache, err = cache.NewRedisDB(&a.cfg.Repository.Redis)
	if err != nil {
		return err
	}

	//redis
	redisRepo := cache.NewRedisRepo(a.cache, a.cfg.Repository.Redis.Ttl)

	// mongo
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1}, // index key ascending
		Options: options.Index().SetUnique(true),
	}

	bookCollection := a.db.Database(a.cfg.Repository.Mongo.DBName).Collection(a.cfg.Repository.Mongo.BooksCollection)
	userCollection := a.db.Database(a.cfg.Repository.Mongo.DBName).Collection(a.cfg.Repository.Mongo.UsersCollection)

	_, err = userCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}

	mongoRepo := db.NewRepoMongo(userCollection, bookCollection)

	//jwt and hasher
	tokenManager, err := auth.NewManager(a.cfg.Auth.JWT.SigningKey)
	if err != nil {
		return err
	}
	hasher := hash.NewSHA1Hasher(a.cfg.Auth.PasswordSalt)

	// services
	userService := user_service.NewUserService(mongoRepo, hasher, tokenManager, a.cfg.Auth.JWT.AccessTokenTTL, a.cfg.Auth.JWT.RefreshTokenTTL)
	bookService := book_service.NewBookService(mongoRepo)

	a.router.Get("/swagger/*", httpSwagger.WrapHandler)
	responder := http2.NewResponder(a.logger)

	v1.SetHandler(a.router, responder, a.logger, userService, bookService, redisRepo, tokenManager)

	return nil
}

func (a *App) Run(ctx context.Context) {
	var err error
	defer a.closeConnections()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:        a.router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    a.cfg.App.RTO,
		WriteTimeout:   a.cfg.App.WTO,
	}
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			a.logger.Info(err.Error())
			cancel()
			return
		}
	}()
	a.logger.Info("started", zap.Int("port", a.cfg.App.Port))
	<-ctx.Done()
	a.logger.Info("shutting down server ...\n")

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		a.logger.Info("server forced shutdown ...")
	}

	a.logger.Info("server exiting ...")
}

func (a *App) closeConnections() {
	defer a.logger.Sync()
	if err := a.db.Disconnect(context.Background()); err != nil {
		a.logger.Error(err.Error())
	}
	if err := a.cache.Close(); err != nil {
		a.logger.Error(err.Error())
	}
	// если будет бд и т.д, закрывать коннект здесь
}
