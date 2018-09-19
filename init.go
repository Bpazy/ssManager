package main

import (
	"database/sql"
	"flag"
	"github.com/Bpazy/ssManager/ss"
	"github.com/Bpazy/ssManager/util"
	"strings"
)

func init() {
	dbPath = flag.String("dbPath", "temp.db", "sqlite数据库文件路径")
	port = flag.String("port", ":8082", "server port")
	version := flag.String("type", "go", "shadowsocks type. such python or go")
	configFilename = flag.String("filename", "config.json", "shadowsocks config file name")
	flag.Parse()

	db2, err := sql.Open("sqlite3", *dbPath)
	util.ShouldPanic(err)
	db = db2

	sc = createSsClient(*version)

	if !tableExists("s_ports") {
		createTable("s_ports")
	}
	if !tableExists("s_user") && !tableExists("s_user_password") {
		createTable("s_user")
		createTable("s_user_password")
		SaveUser(createAdminUser())
	}
}

func createSsClient(version string) ss.Client {
	switch version {
	case "go":
		return &ss.GoClient{Filename: *configFilename}
	}
	panic("not support shadowsocks version")
}

func createAdminUser() (*User, string) {
	return &User{"1", "admin", "admin", "admin@admin.com"}, "admin"
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
