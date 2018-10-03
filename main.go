package main

import (
	"database/sql"
	"github.com/Bpazy/ssManager/cookie"
	"github.com/Bpazy/ssManager/iptables"
	"github.com/Bpazy/ssManager/result"
	"github.com/Bpazy/ssManager/ss"
	"github.com/Bpazy/ssManager/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	tableInitPath = "res/init_{}.sql"
)

var (
	db             *sql.DB
	port           *string
	dbPath         *string
	sc             ss.Client
	configFilename *string
)

func main() {

	r := gin.Default()
	r.Use(errorMiddleware(), authMiddleware())
	api := r.Group("/api")
	{
		api.POST("/login", loginHandler())
		api.GET("/echo", echo())
	}

	usageApi := api.Group("/usage")
	{
		usageApi.GET("/list", listHandler())
		usageApi.POST("/save", saveHandler())
		usageApi.POST("/edit", editHandler())
		usageApi.GET("/delete/:port", deleteHandler())
		usageApi.GET("/reset/:port", resetHandler())
	}

	ssApi := api.Group("/ss")
	{
		ssApi.POST("/addPortPassword", addPortPasswordHandler())
		ssApi.GET("/queryPortPasswords", queryPortsHandler())
		ssApi.GET("/deletePort/:port", deletePortHandler())
		ssApi.POST("/restart", restartHandler())
	}

	r.Run(*port)
}

func monitor() {
	ports := QueryPorts()
	intPorts := make([]int, len(ports))

	for _, p := range ports {
		intPorts = append(intPorts, p.Port)
	}

	usageMap := iptables.GetSptUsageMap(intPorts)
}

func echo() gin.HandlerFunc {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return func(c *gin.Context) {
		con, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}

		con.ReadMessage()

		go func() {
			for {
				err = con.WriteMessage(websocket.TextMessage, []byte("la la la"))
				if err != nil {
					panic(err)
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
}

func queryPortsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, sc.QueryPortPasswords())
	}
}

func deletePortHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		port := c.Param("port")
		sPort, err := strconv.Atoi(port)
		if err != nil {
			c.JSON(http.StatusOK, result.Fail("number port required", ""))
		}
		DeletePort(sPort)
		c.JSON(http.StatusOK, result.Ok("delete success", ""))
	}
}

func restartHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func addPortPasswordHandler() gin.HandlerFunc {
	type addPortPasswordRequest struct {
		Port     string `json:"port"`
		Password string `json:"password"`
	}
	return func(c *gin.Context) {
		r := addPortPasswordRequest{}
		if ok := util.BindJson(c, &r); !ok {
			c.JSON(http.StatusOK, result.Fail("add port password failed", ""))
			return
		}
		AddPortPassword(r.Port, r.Password)
		c.JSON(http.StatusOK, result.Ok("add port password success", ""))
	}
}

func loginHandler() gin.HandlerFunc {
	type loginRequest = struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(c *gin.Context) {
		l := loginRequest{}
		if ok := util.BindJson(c, &l); !ok {
			return
		}
		user := FindUserByAuth(l.Username, l.Password)
		if user == nil {
			c.JSON(http.StatusOK, result.Fail("username or password is incorrect", ""))
			return
		}
		cookie.SaveUserId(c, user.UserId)
		c.JSON(http.StatusOK, result.Ok("login success", user))
	}
}

var whiteList = []string{"/api/login", "/api/echo"}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if util.ContainsString(whiteList, c.Request.RequestURI) {
			return
		}
		userId, err := cookie.GetSavedUserId(c)
		if err != nil {
			c.JSON(419, result.Fail("login needed", err.Error()))
			c.Abort()
			return
		}
		user := FindUser(userId)
		if user == nil {
			c.JSON(419, result.Fail("login needed", err.Error()))
			cookie.Clear(c)
			c.Abort()
		}
	}
}

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(-1, result.Fail(c.Errors.Last().Error(), ""))
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
		portStructs := QueryPortsWithUsage()
		sort.Sort(sort.Reverse(PortSorter(portStructs)))

		c.JSON(http.StatusOK, portStructs)
	}
}
