package sqlite

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sqlx.DB

	DSN string

	// useful for testing & setting a uniform timestamp
	Now func() time.Time
}

// Open opens the database connection.
func (db *DB) Open() error {
	if db.Now == nil {
		db.Now = time.Now
	}

	db.db = sqlx.MustConnect("sqlite3", db.DSN)

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.db.Close()
}
