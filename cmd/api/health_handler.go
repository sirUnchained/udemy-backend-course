package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	if err := writeJSONError(w, 404, "404"); err != nil {
		app.internalServerError(w, r, err)
	}
}
