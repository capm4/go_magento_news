package database

import (
	"context"
	"database/sql"
	"magento/bot/pkg/model"

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

// get all websites
func (p *PostgresWebsites) GetAll(ctx context.Context) (*sql.Rows, error) {
	rows, err := getAll(*p.db, tableNameWebsite, ctx)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// get document by id
func (p *PostgresWebsites) GetById(id int64, ctx context.Context) (*sql.Row, error) {
	whereStm := []PostgresWhereParam{}
	whereStm = append(whereStm, PostgresWhereParam{Column: "id", Type: "eq", Value: id})
	query, _, err := createSelectWhereStm(tableNameWebsite, whereStm).ToSQL()
	if err != nil {
		return nil, err
	}
	row := p.db.client.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}

// update website
func (p *PostgresWebsites) Update(website model.Website, ctx context.Context) (int64, error) {
	rowsAffected, err := update(*p.db, website, website.Id, tableNameWebsite, ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

// update document with value
func (p *PostgresWebsites) UpdateById(id int64, url, selector, attribute, last_url string, ctx context.Context) (int64, error) {
	query := createUpdateQuery(tableNameWebsite, "url = $1, selector = $2, attribute = $3, last_url = $4 WHERE id = $5")
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
func (p *PostgresWebsites) DeleteById(id int64, ctx context.Context) (int64, error) {
	rowsAffected, err := deleteById(*p.db, id, tableNameWebsite, ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

// create document by id
func (p *PostgresWebsites) Insert(website model.Website, ctx context.Context) (int64, error) {
	t := WebsiteToDb(&website)
	query, _, err := createInsertStm(tableNameWebsite, t).Returning("id").ToSQL()
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
	row.Scan(&website.Id)
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return website.Id, nil
}

func (p *PostgresWebsites) IsExistById(id int64, ctx context.Context) (bool, error) {
	exists, err := isExistById(*p.db, id, ctx, tableNameWebsite)
	if err != nil {
		return false, err
	}
	return exists, nil
}
