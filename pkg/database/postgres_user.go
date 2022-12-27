package database

import (
	"context"
	"database/sql"
	"fmt"
	"magento/bot/pkg/model"
	"time"

	_ "github.com/lib/pq"
)

const (
	tableNameUser = "users"
)

type PostgresUser struct {
	db *PostgresDB
}

type DbUser struct {
	Name      string     `db:"name"`
	Login     string     `db:"login"`
	Password  string     `db:"password"`
	IsActive  bool       `db:"is_active" default:"true"`
	UserRole  string     `db:"role"`
	CreatedAt *time.Time `db:"-"`
	UpdatedAt *time.Time `db:"updated_at"`
}

//create databases connection
func NewPostgresUserDB(db *PostgresDB) (PostgressUserInterface, error) {
	return &PostgresUser{db: db}, nil
}

func UserToDb(u *model.User) DbUser {
	return DbUser{
		Name:      u.Name,
		Login:     u.Login,
		Password:  u.Password,
		IsActive:  u.IsActive,
		UserRole:  u.UserRole,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// get all users
func (p *PostgresUser) GetAll(ctx context.Context) (*sql.Rows, error) {
	rows, err := getAll(*p.db, tableNameUser, ctx)
	if err != nil {
		return nil, err
	}
	return rows, nil
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
	rowsAffected, err := update(*p.db, user, user.Id, tableNameUser, ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

// create user
func (p *PostgresUser) Insert(user model.User, ctx context.Context) (int64, error) {
	u := UserToDb(&user)
	query, _, err := createInsertStm(tableNameUser, u).Returning("id").ToSQL()
	if err != nil {
		return 0, err
	}
	tx, err := p.db.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx, query)
	row.Err()
	if row.Err() != nil {
		return 0, row.Err()
	}
	row.Scan(&user.Id)
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (p *PostgresUser) IsExistByLogin(login string, ctx context.Context) (bool, error) {
	whereStm := []PostgresWhereParam{}
	whereStm = append(whereStm, PostgresWhereParam{Column: "login", Type: "eq", Value: login})
	query, _, err := createSelectWhereStm(tableNameUser, whereStm).ToSQL()
	if err != nil {
		return false, err
	}
	var exists bool
	q := fmt.Sprintf("SELECT exists (%s)", query)
	p.db.client.QueryRowContext(ctx, q).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

// delete user by id
func (p *PostgresUser) DeleteById(id int64, ctx context.Context) (int64, error) {
	rowsAffected, err := deleteById(*p.db, id, tableNameUser, ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

func (p *PostgresUser) IsExistById(id int64, ctx context.Context) (bool, error) {
	exists, err := isExistById(*p.db, id, ctx, tableNameUser)
	if err != nil {
		return false, err
	}
	return exists, nil
}
