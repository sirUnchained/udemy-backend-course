package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirUnchained/udemy-backend-course/internal/env"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
	cfg := config{addr: env.GetString("ADDR", ":8000")}
	app := &application{config: cfg}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
