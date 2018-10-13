package db

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Ins    *sql.DB
	dbPath = flag.String("dbPath", "temp.db", "sqlite数据库文件路径")
)

func init() {
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		panic(err)
	}
	Ins = db
}
