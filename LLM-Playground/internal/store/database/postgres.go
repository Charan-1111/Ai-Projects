package database

import (
	"context"
	"fmt"
	"os"
	"sync"

	"llm-playground/internal/logging"
	"llm-playground/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DataBaseStore struct {
	Db      *pgxpool.Pool
	Queries models.Queries
	Log     *logging.Log
	once    sync.Once
}

func (db *DataBaseStore) InitializeDatabaseStore(ctx context.Context, queries models.Queries, log *logging.Log) error {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	db.Queries = queries
	db.Log = log

	var err error

	db.once.Do(func() {
		dsn := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbName
		db.Db, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return
		}
	})

	if err != nil {
		return fmt.Errorf("Error initializing database connection pool : %w", err)
	}

	return nil
}

func (db *DataBaseStore) PingDb(ctx context.Context) error {
	return db.Db.Ping(ctx)
}

func (db *DataBaseStore) Close() {
	if db.Db != nil {
		db.Db.Close()
	}
}
