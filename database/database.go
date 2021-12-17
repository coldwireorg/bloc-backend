package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

// Function to connect to the database
func Connect() error {
	var err error
	DB, err = pgxpool.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return err
	}

	return nil
}
