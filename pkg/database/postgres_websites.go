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
	tableNameWebsite = "websites"
)

type PostgresWebsites struct {
	db *PostgresDB
}

type PostgreWebsite struct {
	Url       string `db:"url"`
	Selector  string `db:"selector"`
	Attribute string `db:"attribute"`
	LastUrl   string `db:"last_url"`
}

func WebsiteToDb(w *model.Website) PostgreWebsite {
	return PostgreWebsite{Url: w.Url, Selector: w.Selector, Attribute: w.Attribute, LastUrl: w.LastUrl}
}

//create databases connection
func NewPostgresWebsiteDB(db *PostgresDB) (PostgressWebsitesInterface, error) {
	return &PostgresWebsites{db: db}, nil
}

// get all documents
func (p *PostgresWebsites) GetAll() (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	query, _, err := createSelectStm(tableNameWebsite).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := p.db.client.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// get document by id
func (p *PostgresWebsites) GetById(id int64) (*sql.Row, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	whereStm := []PostgresWhereParam{}
	whereStm = append(whereStm, PostgresWhereParam{Column: "id", Type: "eq", Value: id})
	query, _, err := createSelectWhereStm(tableNameWebsite, whereStm).ToSQL()
	if err != nil {
		return nil, err
	}
	row := p.db.client.QueryRowContext(ctx, query)

	return row, nil
}

// update website
func (p *PostgresWebsites) Update(website model.Website) (int64, error) {
	websites := append([]model.Website{}, website)
	query, _, err := createUpdateStm(tableNameWebsite, websites).ToSQL()
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

// update document with value
func (p *PostgresWebsites) UpdateById(id int64, url, selector, attribute, last_url string) (int64, error) {
	query := createUpdateQuery(tableNameWebsite, "url = $1, selector = $2, attribute = $3, last_url = $4 WHERE id = $5")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, err := p.db.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, query, url, selector, attribute, last_url, id)
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

// delete website by id
func (p *PostgresWebsites) DeleteById(id int64) (int64, error) {
	query, _, err := createRemoveStm(tableNameWebsite, "id", id).ToSQL()
	if err != nil {
		return 0, err
	}
	fmt.Println(query)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

// create document by id
func (p *PostgresWebsites) Insert(website model.Website) (int64, error) {
	t := WebsiteToDb(&website)
	query, _, err := createInsertStm(tableNameWebsite, t).Returning("id").ToSQL()
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
	row.Scan(&website.Id)
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return website.Id, nil
}
