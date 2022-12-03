package database

import (
	"database/sql"
	"errors"
	"fmt"
	"magento/bot/pkg/model"

	"github.com/doug-martin/goqu/v9"
)

type PostgresDB struct {
	client *sql.DB
}

type GenerigParamsModel interface {
	model.Website | model.Config | model.SlackBot | model.User | PostgreWebsite
}

type GenerigParams interface {
	int | int8 | int16 | int32 | int64 | string
}

type PostgresWhereParam struct {
	Column string
	Type   string
	Value  interface{}
}

//create databases connection
func NewPostgresDB(host, port, user, password, dbname string) (*PostgresDB, error) {
	if host == "" || port == "" {
		return nil, errors.New("mongo DB host and port is required")
	}
	psqlconn := fmt.Sprintf("host= %s port = %s user = %s password = %s dbname = %s sslmode=disable", host, port, user, password, dbname)
	client, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{client: client}, nil
}

//this function create main select (SELECT * FROM table)
//parament will be add after main select
//if add to function WHERE ID = ?
//result will be SELECT * FROM table WHERE ID = ?
func createSelect(tableName, query string) string {

	return fmt.Sprintf("SELECT * FROM %s %s", tableName, query)
}

func createSelectStm(table string) *goqu.SelectDataset {
	return goqu.From(table)
}

func createSelectWhereStm(table string, params []PostgresWhereParam) *goqu.SelectDataset {
	stmExp := goqu.Ex{}
	for _, param := range params {
		addExToStm(stmExp, param)
	}
	sql := goqu.From(table).Where(stmExp)
	return sql
}

func addJoin(stm *goqu.SelectDataset, joinTable, fKey, joinKey string) *goqu.SelectDataset {
	return stm.Join(goqu.T(joinTable), goqu.On(goqu.Ex{fKey: goqu.I(joinKey)}))
}

//table is table whitch to update
//updateParams is slice of map[string]string
//example
// maps := make([]map[string]string)
// m := make(map[string]string)
// m["url"] = "test.com"
// m["selector"] = "href"
// m1 := make(map[string]string)
// m1["url"] = "test1.com"
// m1["selector"] = ".a"
// maps[1] = m
// maps[2] = m1
func createUpdateStm[G GenerigParamsModel](table string, updateParams []G) *goqu.UpdateDataset {
	return goqu.Update(table).Set(updateParams)
}

func createInsertStm[G GenerigParamsModel](table string, insertParams G) *goqu.InsertDataset {
	return goqu.Insert(table).Rows(insertParams)
}

//create remove query
//where table is table name where to remove
//colums and value is where statment like remove where colum eq value
func createRemoveStm[ValueParam GenerigParams](table, column string, value ValueParam) *goqu.DeleteDataset {
	ds := goqu.From(table)
	return ds.Where(goqu.C(column).Eq(value)).Delete()
}

func addExToStm(stm goqu.Ex, param PostgresWhereParam) {
	switch pType := param.Type; pType {
	case "eq":
		stm[param.Column] = param.Value
	case "neq":
		stm[param.Column] = goqu.Op{"neq": param.Value}
	case "gt":
		stm[param.Column] = goqu.Op{"gt": param.Value}
	case "lt":
		stm[param.Column] = goqu.Op{"lt": param.Value}
	default:
		stm[param.Column] = param.Value
	}
}

//create Update script
func createUpdateQuery(tableName, query string) string {

	return fmt.Sprintf("UPDATE %s SET %s", tableName, query)
}

//create Delete script
func createDeleteQuery(tableName, query string) string {

	return fmt.Sprintf("DELETE FROM %s %s", tableName, query)
}

//create Insert script
func createInsertQuery(tableName, query string) string {

	return fmt.Sprintf("INSERT INTO %s %s", tableName, query)
}
