package db

import (
	"flag"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultDBFile = "/data/message.db"
)

var (
	dbFile = ""
)

func init() {
	flag.StringVar(&dbFile, "db", defaultDBFile, "default message db file")
}

type DbOperator struct {
	DB *sqlx.DB
}

func NewDbOperator() *DbOperator {
	source := fmt.Sprintf("file:%s?cache=shared", dbFile)
	db := sqlx.MustOpen("sqlite3", source)
	db.SetMaxOpenConns(1)

	return &DbOperator{DB: db}
}

func (db *DbOperator) Close() error {
	return db.DB.Close()
}
