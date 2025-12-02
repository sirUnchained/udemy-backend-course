package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
