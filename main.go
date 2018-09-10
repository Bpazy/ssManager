package main

import (
	"github.com/Bpazy/ssManager/result"
	"github.com/Bpazy/ssManager/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
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
		group.POST("/save", saveHandler())
	}
	r.Run(":7777")
}

func saveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := Port{}
		if ok := util.BindJson(c, &p); !ok {
			return
		}
		SavePort(&p)
		c.JSON(http.StatusOK, result.Ok("save success", &p))
	}
}

func listHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		portStructs := QueryPorts()
		c.JSON(http.StatusOK, portStructs)
	}
}
