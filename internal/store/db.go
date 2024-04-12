package store

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/khalil-farashiani/golim/internal/store/role"
	"log"
	"strings"
)

type Store struct {
	db    *role.Queries
	cache *Cache
}

//go:embed schema.sql
var ddl string

func InitDB(ctx context.Context) *sql.DB {
	db, err := sql.Open("sqlite3", "golim.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatal(err)
	}
	return db
}
