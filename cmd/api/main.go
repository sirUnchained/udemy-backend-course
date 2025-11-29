package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirUnchained/udemy-backend-course/internal/db"
	"github.com/sirUnchained/udemy-backend-course/internal/env"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:strongpassword@localhost:5432/postgres?sslmode=disable"), // sslmode is disable because we are running locally
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),                                                                    // Maximum number of open connections to the database
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),                                                                    // Maximum number of idle connections in the pool
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),                                                               // Maximum time a connection can remain idle before being closed
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		panic(err)
	}
	store := store.NewPostgresStorage(db)
	// HOLLY SHI*T! i forgot to close database!!
	defer db.Close()
	log.Println("database connected.")

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
