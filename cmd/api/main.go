package main

import (
	"log"

	"github.com/sirUnchained/udemy-backend-course/internal/env"
)

func main() {
	cfg := config{addr: env.GetString("ADDR", ":8000")}
	app := &application{config: cfg}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
