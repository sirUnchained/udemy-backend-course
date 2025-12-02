package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// recovers from panics and returns 500 error
	r.Use(middleware.Recoverer)
	// a simple logger for HTTP requests
	r.Use(middleware.Logger)
	// sets real IP from X-Real-IP or X-Forwarded-For headers
	r.Use(middleware.RealIP)
	// adds a unique request ID to each request
	r.Use(middleware.RequestID)
	// sets timeout for requests to prevent hanging connections
	r.Use((middleware.Timeout(time.Second * 60)))

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// health check endpoint
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postid}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostByIdHandler)
				r.Delete("/", app.deletePostByIdHandler)
				r.Patch("/", app.updatePostByIdHandler)
			})
		})

		r.Route("/comments", func(r chi.Router) {
			r.Post("/post/{postid}", app.createCommentHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userid}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
		})
	})

	// chi.Mux implements http.Handler interface, we have no error if we return it
	return r
}

func (app *application) run(mux http.Handler) error {
	// Configure HTTP server settings
	srv := &http.Server{
		Addr:         app.config.addr,  // Server address to listen on
		Handler:      mux,              // request router/handler
		WriteTimeout: time.Second * 30, // maximum duration before timing out writes
		ReadTimeout:  time.Second * 10, // maximum duration before timing out reads
		IdleTimeout:  time.Minute * 1,  // maximum idle connection timeout
	}

	// Log server start information
	log.Printf("server listening on %s", app.config.addr)

	// Start the HTTP server
	return srv.ListenAndServe()
}
