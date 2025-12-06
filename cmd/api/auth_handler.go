package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,max=100,min=8"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler		godoc
//
//	@Summary		register user
//	@Description	i dont know! just register a fucking user!
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	UserWithToken			true	"User credentials"
//	@Success		200	{object}	store.User	"user registered"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	newUser := &store.User{
		UserName: payload.Username,
		Email:    payload.Email,
	}

	// hashe first then set hashed password
	if err := newUser.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	// hash token for storage but keep the plain token for email
	hashed := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hashed[:])

	ctx := r.Context()
	if err := app.store.Users.CreateAndInvite(ctx, newUser, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicatedEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicatedUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := &UserWithToken{
		User:  newUser,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusOK, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// activateUserHandler		godoc
//
//	@Summary		activate/register user
//	@Description	i dont know! just activate/register a fucking user!
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path	string			true	"invitation credentials"
//	@Success		200	{object}	store.User	"user activated"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [post]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.ErrorNoRow:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
