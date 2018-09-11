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
		group.GET("/delete/:port", deleteHandler())
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
		DeletePort(port)
		c.JSON(http.StatusOK, result.Ok("删除成功", ""))
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
		SaveIptables(p.Port)
		c.JSON(http.StatusOK, result.Ok("save success", &p))
	}
}

func listHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		portStructs := QueryPorts()
		c.JSON(http.StatusOK, portStructs)
	}
}
