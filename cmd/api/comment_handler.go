package main

import (
	"net/http"

	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

type commentPayload struct {
	UserID  int64  `json:"user_id" validate:"required,numeric"`
	PostID  int64  `json:"post_id" validate:"required,numeric"`
	Content string `json:"content" validate:"required,max=512"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload commentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		UserID:  payload.UserID,
		PostID:  payload.PostID,
		Content: payload.Content,
	}

	ctx := r.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, "comment added."); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
