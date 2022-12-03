package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	tableNameConfig = "config"
)

type PostgresConfig struct {
	db *PostgresDB
}

//create databases connection
func NewPostgresConfigDB(db *PostgresDB) (PostgressConfigInterface, error) {
	return &PostgresConfig{db: db}, nil
}

// get document by id
func (p *PostgresConfig) GetByPath(path string) *sql.Row {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	row := p.db.client.QueryRowContext(ctx, createSelect(tableNameConfig, "WHERE path = $1"), path)

	return row
}

// update document with value
func (p *PostgresConfig) UpdateById(id int64, path, value string) (int64, error) {
	query := createUpdateQuery(tableNameConfig, "path = $1, value = $2")
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	tx, err := p.db.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, query, path, value)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

// delete document by id
func (p *PostgresConfig) DeleteById(id int64) (int64, error) {
	query := createDeleteQuery(tableNameConfig, "WHERE id = $1")
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	tx, err := p.db.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

// create config
func (p *PostgresConfig) Insert(path, value string) (int64, error) {
	query := createInsertQuery(tableNameConfig, "(path, value) VALUES ($1, $2)")
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	tx, err := p.db.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, query, path, value)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
