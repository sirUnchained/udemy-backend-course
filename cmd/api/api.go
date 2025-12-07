package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirUnchained/udemy-backend-course/docs"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	apiURL string
	mail   mailConfig
	auth   authConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	exp time.Duration
}

type authConfig struct {
	basic basicConfig
}

type basicConfig struct {
	user string
	pass string
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
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

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
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userid}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})

		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})
	})

	// chi.Mux implements http.Handler interface, we have no error if we return it
	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = VERSION
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	// Configure HTTP server settings
	srv := &http.Server{
		Addr:         app.config.addr,  // Server address to listen on
		Handler:      mux,              // request router/handler
		WriteTimeout: time.Second * 30, // maximum duration before timing out writes
		ReadTimeout:  time.Second * 10, // maximum duration before timing out reads
		IdleTimeout:  time.Minute * 1,  // maximum idle connection timeout
	}

	// Log server start information
	app.logger.Infoln("server started", "addr:", app.config.addr)

	// Start the HTTP server
	return srv.ListenAndServe()
}
