package database

import (
	"context"
	"database/sql"
	"fmt"
	"magento/bot/pkg/model"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

const (
	tableNameSlack        = "slack_bot"
	tableNameSlackWebsite = "slack_bot_websites"
)

type PostgresSlackBot struct {
	db *PostgresDB
}

type PostgreSlackBot struct {
	Name      string `db:"name"`
	Token     string `db:"token"`
	ChannelId string `db:"channel_id"`
	CronEvery int64  `db:"cron_every"`
}

type PostgreSlackBotWithTime struct {
	Name        string     `db:"name"`
	Token       string     `db:"token"`
	ChannelId   string     `db:"channel_id"`
	CronEvery   int64      `db:"cron_every"`
	LastCronRun *time.Time `db:"last_cron_run"`
}

type PostgresSlackWebiste struct {
	SlackId   int64 `db:"slack_id"`
	WebsiteId int64 `db:"website_id"`
}
type PostgresSlackWebisteReturn struct {
	Id int64 `db:"id"`
}

func SlackBockToDb(s *model.SlackBot) PostgreSlackBot {
	return PostgreSlackBot{Name: s.Name, Token: s.Token, ChannelId: s.ChannelId, CronEvery: s.CronEvery}
}

func SlackWebsiteDb(slackId, websiteId int64) PostgresSlackWebiste {
	return PostgresSlackWebiste{SlackId: slackId, WebsiteId: websiteId}
}

func SlackBockToDbWithTime(s *model.SlackBot) PostgreSlackBotWithTime {
	return PostgreSlackBotWithTime{
		Name:        s.Name,
		Token:       s.Token,
		ChannelId:   s.ChannelId,
		CronEvery:   s.CronEvery,
		LastCronRun: s.LastCronRun,
	}
}

//create databases connection
func NewPostgresSlackBotDB(db *PostgresDB) (PostgressSlackInterface, error) {
	return &PostgresSlackBot{db: db}, nil
}

// get all slackBot
func (p *PostgresSlackBot) GetAll(ctx context.Context) (*sql.Rows, error) {
	rows, err := getAll(*p.db, tableNameSlack, ctx)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// get slack bot by id
func (p *PostgresSlackBot) GetById(id int64, ctx context.Context) (*sql.Row, error) {
	whereStm := []PostgresWhereParam{}
	whereStm = append(whereStm, PostgresWhereParam{Column: "id", Type: "eq", Value: id})
	query, _, err := createSelectWhereStm(tableNameSlack, whereStm).ToSQL()
	if err != nil {
		return nil, err
	}
	row := p.db.client.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}

// create slack bot
func (p *PostgresSlackBot) Insert(slack model.SlackBot, ctx context.Context) (int64, error) {
	t := SlackBockToDbWithTime(&slack)
	query, _, err := createInsertStm(tableNameSlack, t).Returning("id").ToSQL()
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
	row.Scan(&slack.Id)
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return slack.Id, nil
}

func (p *PostgresSlackBot) IsExistById(id int64, ctx context.Context) (bool, error) {
	exists, err := isExistById(*p.db, id, ctx, tableNameSlack)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// update slackbot
func (p *PostgresSlackBot) Update(slackbot model.SlackBot, ctx context.Context) (int64, error) {
	rowsAffected, err := update(*p.db, slackbot, slackbot.Id, tableNameSlack, ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

// delete slack bot by id
func (p *PostgresSlackBot) DeleteById(id int64, ctx context.Context) (int64, error) {
	rowsAffected, err := deleteById(*p.db, id, tableNameSlack, ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

func (p *PostgresSlackBot) InsertWebsiteToSlack(slackId, websiteId int64, ctx context.Context) (int64, error) {
	t := SlackWebsiteDb(slackId, websiteId)
	query, _, err := createInsertStm(tableNameSlackWebsite, t).Returning("id").ToSQL()
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
	returnData := PostgresSlackWebisteReturn{}
	row.Scan(&returnData.Id)
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return returnData.Id, nil
}

func (p *PostgresSlackBot) IsExistWebsiteInSlack(slackId, websiteId int64, ctx context.Context) (bool, error) {
	whereStm := []PostgresWhereParam{}
	whereStm = append(whereStm, PostgresWhereParam{Column: "slack_id", Type: "eq", Value: slackId})
	whereStm = append(whereStm, PostgresWhereParam{Column: "website_id", Type: "eq", Value: websiteId})
	query, _, err := createSelectWhereStm(tableNameSlackWebsite, whereStm).ToSQL()
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

func (p *PostgresSlackBot) DeleteWebsiteFromSlackById(slackId, websiteId int64, ctx context.Context) (int64, error) {
	whereStm := []PostgresWhereParam{}
	whereStm = append(whereStm, PostgresWhereParam{Column: "slack_id", Type: "eq", Value: slackId})
	whereStm = append(whereStm, PostgresWhereParam{Column: "website_id", Type: "eq", Value: websiteId})
	rowsAffected, err := deleteByMultimpleColumn(*p.db, whereStm, tableNameSlackWebsite, ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, err
}

func (p *PostgresSlackBot) GetAllWebsiteBySlackId(id int64, ctx context.Context) (*sql.Rows, error) {
	sQ := goqu.From(tableNameWebsite).Select("websites.id", "url", "selector", "attribute", "last_url")
	sQ = sQ.LeftJoin(goqu.I(tableNameSlackWebsite), goqu.On(goqu.Ex{"slack_bot_websites.website_id": goqu.I("websites.id")}))
	sQ = sQ.Where(goqu.C("slack_id").Eq(id))
	query, _, err := sQ.ToSQL()
	if err != nil {
		return nil, err
	}
	fmt.Println(query)
	rows, err := p.db.client.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
