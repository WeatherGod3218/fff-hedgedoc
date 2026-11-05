package database

import (
	"context"
	"fff-hedgedoc/internal/logging"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool
var DatabaseInitialized bool = false

func InitDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	addr := os.Getenv("POSTGRES_ADDR")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	connURL := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(user, password),
		Host:   addr,
		Path:   "/" + user,
	}

	newDB, err := pgxpool.New(
		ctx,
		connURL.String(),
	)

	if err != nil {
		return err
	}

	err = newDB.Ping(ctx)
	if err != nil {
		return err
	}
	db = newDB

	DatabaseInitialized = true

	logging.Logger.Info("Database has been initalized!")

	return nil
}
