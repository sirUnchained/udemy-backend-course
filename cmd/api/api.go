package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// it is a simple logger
	r.Use(middleware.Logger)
	// it is going to recover from panics
	r.Use(middleware.Recoverer)
	// it is going to set a real ip from X-Real-IP or X-Forwarded-For headers
	r.Use(middleware.RealIP)
	// it is going to set a request id for each request
	r.Use(middleware.RequestID)
	// it is a timeout for requests
	r.Use((middleware.Timeout(time.Second * 60)))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	// chi.Mux impeliments http.Handler, so there is no problem for return it
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute * 1,
	}

	log.Printf("server listen on %s", app.config.addr)

	return srv.ListenAndServe()
}
