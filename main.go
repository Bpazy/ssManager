package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

const (
	usage = "iptables -L -nvx | grep spt:{} | awk '{print $2}'"
)

var db *sqlx.DB

func init() {
	db2, err := sqlx.Open("sqlite3", ":memory:")
	checkErr(err)
	db = db2

	_, err = db.Exec("CREATE table s_ports (port int);")
	checkErr(err)
	_, err = db.Exec("insert into s_ports (port) values (6666);")
	checkErr(err)
}

func main() {
	r := gin.Default()
	group := r.Group("/api")
	{
		group.GET("/list", listHandler())
	}
	r.Run(":7777")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func queryPorts() []int {
	rows, err := db.Queryx("SELECT * FROM s_ports;")
	checkErr(err)
	defer rows.Close()
	ports := make([]int, 0)
	for rows.Next() {
		var p int
		rows.Scan(&p)
		ports = append(ports, p)
	}
	return ports
}

func listHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ports := queryPorts()
		portUsage := make(map[int]int, len(ports))
		for _, p := range ports {
			replace := strings.Replace(usage, "{}", strconv.Itoa(p), -1)
			fmt.Println(replace)
			cmd := exec.Command("bash", "-c", replace)
			var out bytes.Buffer
			cmd.Stdout = &out

			err := cmd.Run()
			checkErr(err)
			portUsage[p] = parseInt(strings.TrimSpace(out.String()))
		}
		c.JSON(http.StatusOK, portUsage)
	}
}

func parseInt(i string) int {
	atoi, e := strconv.Atoi(i)
	checkErr(e)
	return atoi
}
