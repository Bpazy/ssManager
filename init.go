package main

import (
	"github.com/Bpazy/ssManager/util"
	"github.com/jmoiron/sqlx"
	"strings"
)

func init() {
	db2, err := sqlx.Open("sqlite3", "dev.db")
	util.ShouldPanic(err)
	db = db2

	if !tableExists("s_ports") {
		createTable("s_ports")
	}
}

func createTable(tableName string) {
	_, err := db.Exec(getSql(strings.Replace(tableInitPath, "{}", tableName, -1)))
	util.ShouldPanic(err)
}

func tableExists(tableName string) bool {
	row := db.QueryRow("select count(*) from sqlite_master where name = ?;", tableName)
	count := 0
	err := row.Scan(&count)
	util.ShouldPanic(err)

	if count != 0 {
		return true
	}
	return false
}

func getSql(filePath string) string {
	b, err := Asset(filePath)
	util.ShouldPanic(err)
	return string(b)
}
