package main

import (
	"flag"
	"github.com/Bpazy/ssManager/db"
	"github.com/Bpazy/ssManager/ss"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func init() {
	flag.Parse()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	sc = createSsClient(*version)

	if !tableExists("s_ports") {
		createTable("s_ports")
	}
	if !tableExists("s_user") && !tableExists("s_user_password") {
		createTable("s_user")
		createTable("s_user_password")
		SaveUser(createAdminUser())
	}
	if !tableExists("s_usage") {
		createTable("s_usage")
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
	_, err := db.Ins.Exec(getSql(strings.Replace(tableInitPath, "{}", tableName, -1)))
	if err != nil {
		panic(err)
	}
}

func tableExists(tableName string) bool {
	row := db.Ins.QueryRow("select count(*) from sqlite_master where name = ?;", tableName)
	count := 0
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}

	if count != 0 {
		return true
	}
	return false
}

func getSql(filePath string) string {
	b, err := Asset(filePath)
	if err != nil {
		panic(err)
	}
	return string(b)
}
