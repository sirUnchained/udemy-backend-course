package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirUnchained/udemy-backend-course/internal/db"
	"github.com/sirUnchained/udemy-backend-course/internal/env"
	"github.com/sirUnchained/udemy-backend-course/internal/seeds"
	"github.com/sirUnchained/udemy-backend-course/internal/store"
	"go.uber.org/zap"
)

const VERSION = "1.0"

//	@title			go-social api
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @SecurityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// Logger configs
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// ENV configs
	err := godotenv.Load()
	if err != nil {
		logger.Fatalln(err)
		os.Exit(-1)
	}

	// Set ENV
	debugMode := env.GetInt("DEBUGMODE", 1)
	cfg := config{
		addr:   env.GetString("ADDR", ":8000"),
		apiURL: env.GetString("EXTERNAL_URL", "127.0.0.1:4000"),
		// database ENV
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:strongpassword@localhost:5432/postgres?sslmode=disable"), // sslmode is disable because we are running locally
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),                                                                    // Maximum number of open connections to the database
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),                                                                    // Maximum number of idle connections in the pool
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),                                                               // Maximum time a connection can remain idle before being closed
		},
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}

	// start database connection
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatalln(err)
	}
	store := store.NewPostgresStorage(db)
	defer db.Close() // HOLLY SHI*T! i forgot to close database!!
	logger.Infoln("database connected.")

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	// seeds
	seeds.Seed(app.store, (debugMode == 1), db)

	mux := app.mount()
	logger.Fatalln(app.run(mux))
}
