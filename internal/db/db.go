package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	// open a new database connection pool
	sqlDB, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	// configure connection pool settings
	sqlDB.SetMaxIdleConns(maxIdleConns) // Set maximum number of idle connections
	sqlDB.SetMaxOpenConns(maxOpenConns) // Set maximum number of open connections

	// parse the max idle time duration string
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set maximum time a connection can remain idle
	sqlDB.SetConnMaxIdleTime(duration)

	// create context with timeout for database ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// verify database connection is alive
	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	// Return the configured database connection pool
	return sqlDB, nil
}
