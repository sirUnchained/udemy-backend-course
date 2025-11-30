package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Fatalf("internal server error: %s path: %s error: %s\n", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "something went wrong with us, we'll fix this as soon as we can.")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Fatalf("bad request error: %s path: %s error: %s\n", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Fatalf("404 error: %s path: %s error: %s\n", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "the record not found.")
}
