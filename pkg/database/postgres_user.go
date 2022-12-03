package database

import (
	"context"
	"database/sql"
	"magento/bot/pkg/model"

	_ "github.com/lib/pq"
)

const (
	tableNameUser = "users"
)

type PostgresUser struct {
	db *PostgresDB
}

//create databases connection
func NewPostgresUserDB(db *PostgresDB) (PostgressUserInterface, error) {
	return &PostgresUser{db: db}, nil
}

// get user by login
func (p *PostgresUser) GetByLogin(login string, ctx context.Context) (*sql.Row, error) {
	whereStm := []PostgresWhereParam{}
	whereStm = append(whereStm, PostgresWhereParam{Column: "login", Type: "eq", Value: login})
	query, _, err := createSelectWhereStm(tableNameUser, whereStm).ToSQL()
	if err != nil {
		return nil, err
	}
	row := p.db.client.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return row, nil
}

// update user
func (p *PostgresUser) Update(user model.User, ctx context.Context) (int64, error) {
	users := append([]model.User{}, user)
	query, _, err := createUpdateStm(tableNameUser, users).ToSQL()
	if err != nil {
		return 0, err
	}
	tx, err := p.db.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, query)
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

// create user
func (p *PostgresUser) Insert(user model.User, ctx context.Context) (int64, error) {
	query, _, err := createInsertStm(tableNameUser, user).ToSQL()
	if err != nil {
		return 0, err
	}
	tx, err := p.db.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, query)
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
