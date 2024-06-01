package db

import (
	"database/sql"
)

// Store provides all functions to run individual queries as well as transactions
type Store interface {
	Querier
}

// SQLStore provides all functions to run individual queries as well as transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
