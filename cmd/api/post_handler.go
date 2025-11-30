package main

import (
	"log"
	"net/http"

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
		writeJSONError(w, http.StatusBadRequest, "where is your json?")
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

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

}
