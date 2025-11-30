package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	// read JSON from request body into post struct
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json data.")
		return
	}

	userid := 1 // TODO: get from auth context
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  int64(userid),
		Tags:    payload.Tags,
	}

	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		log.Fatalln(err)
		writeJSONError(w, http.StatusInternalServerError, "we'll fix this as soon as we can.")
		return
	}

	if err := writeJSON(w, http.StatusOK, "created new post."); err != nil {
		log.Fatalln(err)
		writeJSONError(w, http.StatusInternalServerError, "we'll fix this as soon as we can.")
		return
	}
}

func (app *application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	// we use chi.URLParam to get the value of the "postid" URL parameter
	postID := chi.URLParam(r, "postid")

	// convert the postID string to an int64 in decimal base
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	ctx := r.Context()
	post, err := app.store.Posts.GetById(ctx, int64(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNoRow):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			log.Fatal(err)
			writeJSONError(w, http.StatusInternalServerError, "something went wrong.")
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		return
	}
}
