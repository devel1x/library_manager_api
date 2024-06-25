package userService

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	v1 "template/internal/delivery/http/v1"
	"template/internal/dto"
	"template/internal/entity"
	"template/internal/utils"
	"template/pkg/auth"
	"template/pkg/hash"
	"template/pkg/validator"
	"time"
)

type UserService struct {
	userRepo userRepo

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(userRepo userRepo, hasher hash.PasswordHasher, manager auth.TokenManager, accesTokenTTL time.Duration, refreshTokenTTl time.Duration) *UserService {
	return &UserService{
		userRepo:        userRepo,
		hasher:          hasher,
		tokenManager:    manager,
		accessTokenTTL:  accesTokenTTL,
		refreshTokenTTL: refreshTokenTTl,
	}
}

type userRepo interface {
	Authenticate(ctx context.Context) (interface{}, error)
	CreateUser(ctx context.Context, user *entity.User) (interface{}, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session entity.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
	GetByCredentials(ctx context.Context, username, password string) (entity.User, error)
}

func (s *UserService) Login(ctx context.Context, form *v1.UserLoginForm) (entity.Tokens, error) {
	passwordHash, err := s.hasher.Hash(form.Password)
	if err != nil {
		return entity.Tokens{}, err
	}

	user, err := s.userRepo.GetByCredentials(ctx, form.Username, passwordHash)

	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return entity.Tokens{}, err
		}

		return entity.Tokens{}, err
	}
	if passwordHash != user.HashedPassword {
		return entity.Tokens{}, utils.ErrInvalidCredentials
	}

	return s.createSession(ctx, user.ID)
}
func (s *UserService) SignUp(ctx context.Context, form *v1.UserSignupForm) (interface{}, error) {
	adminKey := os.Getenv("ADMIN_KEY")
	form.CheckField(validator.MinChars(form.Username, 5), &form.UserErrors.Username, "Username must be at least 5 chars long")
	form.CheckField(validator.MaxChars(form.Username, 20), &form.UserErrors.Username, "Username must be max 20 chars long")
	form.CheckField(validator.NotBlank(form.Username), &form.UserErrors.Username, "Username cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 5), &form.UserErrors.Password, "Password must be at least 5 chars long")
	form.CheckField(validator.MaxChars(form.Password, 20), &form.UserErrors.Password, "Password must be max 20 chars long")
	form.CheckField(validator.NotBlank(form.Password), &form.UserErrors.Password, "Password cannot be blank")
	if form.Secret != "" {
		form.CheckField(validator.CheckAdmin(form.Secret, adminKey), &form.UserErrors.Secret, "Secret key is not matching")
	}
	if !form.ValidUser() {
		return 0, utils.InvalidForm
	}
	user := dto.FormToUser(*form)
	var err error
	user.HashedPassword, err = s.hasher.Hash(form.Password)
	if err != nil {
		return entity.Tokens{}, err
	}
	id, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, utils.ErrUserAlreadyExists) {
			return 0, utils.ErrUserAlreadyExists
		}
		return 0, err
	}
	return id, nil
}

func (s *UserService) GetUserByID(id string) (*entity.User, error) {
	return s.userRepo.GetUserByID(context.Background(), id)
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (entity.Tokens, error) {
	user, err := s.userRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return entity.Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UserService) createSession(ctx context.Context, userId primitive.ObjectID) (entity.Tokens, error) {
	var (
		res entity.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := entity.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.userRepo.SetSession(ctx, userId, session)

	return res, err
}
