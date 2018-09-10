package main

import (
	"bytes"
	"fmt"
	"github.com/Bpazy/ssManager/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

const (
	usage         = "iptables -L -nvx | grep spt:{} | awk '{print $2}'"
	tableInitPath = "res/init_{}.sql"
)

var (
	db *sqlx.DB
)

func main() {
	r := gin.Default()
	group := r.Group("/api")
	{
		group.GET("/list", listHandler())
	}
	r.Run(":7777")
}

func queryPorts() []int {
	rows, err := db.Queryx("SELECT * FROM s_ports;")
	util.ShouldPanic(err)
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
			util.ShouldPanic(err)
			portUsage[p] = util.MustParseInt(strings.TrimSpace(out.String()))
		}
		c.JSON(http.StatusOK, portUsage)
	}
}
