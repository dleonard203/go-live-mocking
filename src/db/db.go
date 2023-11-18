package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DB is a wrapper for interactions with a SQL database
type DB struct {
	db         *sql.DB
	timeout    time.Duration
	retryCount int
}

// NewDatabase creates a new DB object
func NewDatabase(dataSourceName string, retryCount int, timeout time.Duration) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{
		db:      db,
		timeout: timeout,
	}, nil
}

func (d *DB) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d.timeout)
}

// retriedQuery will attempt to issue the query up to d.retryCount times
func (d *DB) retriedQuery(query string, params ...any) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	for i := 0; i < d.retryCount; i++ {
		// in-line func call to defer the cancel
		func() {
			ctx, cancel := d.getContext()
			defer cancel()
			rows, err = d.db.QueryContext(ctx, query, params...)
		}()

		if err == nil {
			break
		}
	}
	return rows, err
}

// User represents an application user
type User struct {
	Name  string
	Email string
}

// GetUser gets the user row for the provided id
func (d *DB) GetUser(id int) (*User, error) {
	rows, err := d.retriedQuery("SELECT name, email FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Name, &user.Email); err != nil {
			return nil, err
		}
		return &user, nil
	}

	return nil, fmt.Errorf("no user found for id=%d", id)
}
