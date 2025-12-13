package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

// ALERT: use 'validate' tag, NOT 'validator'!!
type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=250"`
	Content string   `json:"content" validate:"required,max=1024"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   string `json:"title" validate:"required,max=250"`
	Content string `json:"content" validate:"required,max=1024"`
}

// creating postKey new type for using in ctx
type postKey string

const postCtx postKey = "POST"

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	// read JSON from request body into post struct
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := app.getUserFromCtx(r)
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  user.ID,
		Tags:    payload.Tags,
	}

	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "created new post."); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetCommentsByPostId(r.Context(), int64(post.ID))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostByIdHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	ctx := r.Context()
	if err := app.store.Posts.DeleteById(ctx, int64(post.ID)); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "post is now deleted"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) updatePostByIdHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post.Title = payload.Title
	post.Content = payload.Content

	ctx := r.Context()
	if err := app.store.Posts.Update(ctx, post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we use chi.URLParam to get the value of the "postid" URL parameter
		postID := chi.URLParam(r, "postid")

		// convert the postID string to an int64 in decimal base
		id, err := strconv.ParseInt(postID, 10, 64)
		if err != nil {
			// app.jsonResponseError(w, http.StatusBadRequest, "invalid post ID")
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := app.store.Posts.GetById(ctx, int64(id))
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				// app.jsonResponseError(w, http.StatusNotFound, err.Error())
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	// insted of 'postCtx' we could use a normal string BUT this is not good
	// we should create a type insted and we did it as you can see
	post, _ := r.Context().Value(postCtx).(*store.Post)

	return post
}
