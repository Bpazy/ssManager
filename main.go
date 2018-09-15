package main

import (
	"github.com/Bpazy/ssManager/iptables"
	"github.com/Bpazy/ssManager/result"
	"github.com/Bpazy/ssManager/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"sort"
)

const (
	tableInitPath = "res/init_{}.sql"
)

var (
	db     *sqlx.DB
	port   *string
	dbPath *string
)

func main() {
	r := gin.Default()
	r.Use(errorMiddleware())
	group := r.Group("/api")
	{
		group.GET("/list", listHandler())
		group.POST("/save", saveHandler())
		group.POST("/edit", editHandler())
		group.GET("/delete/:port", deleteHandler())
		group.GET("/reset/:port", resetHandler())
	}
	r.Run(*port)
}

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			panic(c.Errors)
		}
	}
}

func deleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		port := c.Param("port")
		DeletePort(util.MustParseInt(port))
		c.JSON(http.StatusOK, result.Ok("删除成功", ""))
	}
}

func resetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		port := c.Param("port")
		ResetPortUsage(util.MustParseInt(port))
		c.JSON(http.StatusOK, result.Ok("重置成功", ""))
	}
}

func saveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := Port{}
		if ok := util.BindJson(c, &p); !ok {
			return
		}
		ok := SavePort(&p)
		if !ok {
			c.JSON(http.StatusOK, result.Fail("save failed", &p))
		}
		iptables.SaveIptables(p.Port)
		c.JSON(http.StatusOK, result.Ok("save success", &p))
	}
}

func editHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := Port{}
		if ok := util.BindJson(c, &p); !ok {
			return
		}
		ok := EditPort(&p)
		if !ok {
			c.JSON(http.StatusOK, result.Fail("save failed", &p))
		}
		c.JSON(http.StatusOK, result.Ok("edit success", &p))
	}
}

func listHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		portStructs := QueryPorts()
		sort.Sort(sort.Reverse(PortSorter(portStructs)))

		c.JSON(http.StatusOK, portStructs)
	}
}
