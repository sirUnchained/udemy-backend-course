package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

type userKey string

const userCtx userKey = "USER"

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromCtx(r)

	// todo: back when u finished auth
	var payload FollowUser
	fmt.Println(payload)
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromCtx(r)

	// todo: back when u finished auth
	var payload FollowUser
	fmt.Println(payload)
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.UnFollow(ctx, followerUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id_str := chi.URLParam(r, "userid")

		id, err := strconv.ParseInt(id_str, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNoRow):
				app.notFoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
