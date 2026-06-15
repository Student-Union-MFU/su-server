// Package config is used for initializing Database and other configs
package config

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), dsn) 

	if err != nil {
		return  nil, err
	}

	return conn, nil
}


