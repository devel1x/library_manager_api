package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"template/internal/entity"
	"template/internal/utils"
	"template/pkg/validator"
)

type UserLoginForm struct {
	Username            string `json:"username" bson:"username"`
	Password            string `json:"password" bson:"password"`
	validator.Validator `json:"-" bson:"-"`
}

// @Summary User Login
// @Description User login endpoint
// @Tags User
// @Accept json
// @Produce json
// @Param login body UserLoginForm true "Login form"
// @Success 200 {object} entity.Tokens
// @Failure 400 {object} entity.UserFormError "Invalid input"
// @Failure 404 {string} user not found "User not found"
// @Failure 500 {string} Internal server error "Internal server error"
// @Router /api/v1/user/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var form UserLoginForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.responder.WithBadRequest(w, http.StatusText(http.StatusBadRequest))
		return
	}

	res, err := h.userService.Login(ctx, &form)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			h.responder.WithBadRequest(w, err.Error())
			return
		}

		h.responder.WithInternalError(w, err.Error())
		return
	}

	h.responder.WithOK(w, res)
}

type UserSignupForm struct {
	Username            string `json:"username" bson:"username"`
	Password            string `json:"password" bson:"password"`
	Secret              string `json:"secret,omitempty" bson:"secret,omitempty"`
	validator.Validator `json:"-" bson:"-"`
}

// @Summary User Signup
// @Description User signup endpoint
// @Tags User
// @Accept json
// @Produce json
// @Param signup body UserSignupForm true "Signup form"
// @Success 201 {string} user successfully created "User created"
// @Failure 400 {object} entity.UserFormError "Invalid input"
// @Failure 409 {string} user already exists "User already exists"
// @Failure 500 {string} Internal server error "Internal server error"
// @Router /api/v1/user/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var form UserSignupForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.responder.WithBadRequest(w, http.StatusText(http.StatusBadRequest))
		return
	}
	id, err := h.userService.SignUp(ctx, &form)
	if err != nil {
		if errors.Is(err, utils.InvalidForm) {
			h.responder.WriteResponse(w, form.UserErrors, http.StatusBadRequest)
			return
		}
		if errors.Is(err, utils.ErrUserAlreadyExists) {
			h.responder.WithBadRequest(w, err.Error())
			return
		}
		h.responder.WithInternalError(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	h.responder.WithCreated(w, id)
}

// @Summary Refresh user token
// @Description Refresh JWT tokens for a user
// @Tags User
// @Accept json
// @Produce json
// @Param refresh body entity.RefreshInput true "Refresh token information"
// @Success 200 {object} entity.Tokens "New JWT Tokens"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/user/auth/refresh [post]
func (h *Handler) userRefresh(w http.ResponseWriter, r *http.Request) {
	var form entity.RefreshInput
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		h.responder.WithBadRequest(w, "invalid input body")
		return
	}

	res, err := h.userService.RefreshTokens(ctx, form.Token)
	if err != nil {
		h.responder.WithInternalError(w, err.Error())
		return
	}

	h.responder.WithOK(w, res)
}
