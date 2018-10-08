package main

import (
	"database/sql"
	"flag"
	"github.com/Bpazy/ssManager/cookie"
	"github.com/Bpazy/ssManager/id"
	"github.com/Bpazy/ssManager/iptables"
	"github.com/Bpazy/ssManager/net"
	"github.com/Bpazy/ssManager/result"
	"github.com/Bpazy/ssManager/ss"
	"github.com/Bpazy/ssManager/util"
	"github.com/Bpazy/ssManager/ws"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rickb777/date"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	tableInitPath          = "res/init_{}.sql"
	tokenEffectiveDuration = 60

	saveTodayUsageInterval = 1 * time.Minute
)

var (
	db        *sql.DB
	sc        ss.Client
	wsTimeout int64 = 60

	wsList     = arraylist.New()
	wsTokenSet = hashset.New()

	dbPath         = flag.String("dbPath", "temp.db", "sqlite数据库文件路径")
	port           = flag.String("port", ":8082", "server port")
	version        = flag.String("type", "go", "shadowsocks type. such python or go")
	configFilename = flag.String("filename", "config.json", "shadowsocks config file name")
	deviceName     = flag.String("i", "eth0", "network interface device name")
)

type WsToken struct {
	Token      string
	CreateTime time.Time
}

func main() {
	go monitor()
	go wsConnMonitor()
	go tokenMonitor()

	net.AddPorts(QueryPorts())
	go net.DetectNet(*deviceName)
	go writeSpeed()
	go saveTodayUsage()

	r := gin.Default()
	r.Use(errorMiddleware(), authMiddleware())

	wsApi := r.Group("/ws")
	{
		wsApi.GET("/echo", echoHandler())
	}

	api := r.Group("/api")
	{
		api.POST("/login", loginHandler())
		api.GET("/token", tokenHandler())
	}

	usageApi := api.Group("/usage")
	{
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

type usage struct {
	Port      int
	Date      date.Date
	DownUsage int
	UpUsage   int
}

func saveTodayUsage() {
	for {
		downSizeMap := net.GetTodayPortDownSize()
		upSizeMap := net.GetTodayPortUpSize()

		for p, dUsage := range downSizeMap {
			row := db.QueryRow("select port, `date`, downUsage, upUsage from s_usage")
			u2 := usage{}
			err := row.Scan(&u2.Port, &u2.Date, &u2.DownUsage, &u2.UpUsage)
			if err == nil {
				db.Exec("insert into s_usage (port, date, downUsage, upUsage) VALUES (?,?,?,?)",
					p, date.Today(), dUsage)
			}
		}

		time.Sleep(saveTodayUsageInterval)
	}
}

type Speed struct {
	DownSpeed float32 `json:"downSpeed"`
	UpSpeed   float32 `json:"upSpeed"`
}

func writeSpeed() {
	for {
		m := make(map[int]Speed)
		for _, p := range QueryPorts() {
			downSpeed, upSpeed := net.GetSpeed(p)
			m[p] = Speed{
				DownSpeed: downSpeed,
				UpSpeed:   upSpeed,
			}
		}

		for _, e := range wsList.Values() {
			c := e.(*ws.Client)
			c.Conn.WriteJSON(ws.H{
				Type: ws.Speed,
				Data: m,
			})
		}

		time.Sleep(1 * time.Second)
	}
}

func tokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		next := id.Next()
		wsTokenSet.Add(WsToken{next, time.Now()})
		c.JSON(http.StatusOK, next)
	}
}

func monitor() {
	for {
		// TODO performance optimize
		portStructs := QueryPortStructsWithUsage()
		for _, e := range wsList.Values() {
			c := e.(*ws.Client)

			sort.Sort(sort.Reverse(PortStructSorter(portStructs)))
			c.Conn.WriteJSON(ws.H{
				Type: ws.Usage,
				Data: portStructs,
			})
		}

		time.Sleep(3 * time.Second)
	}
}

func tokenMonitor() {
	for {
		for _, wsToken := range wsTokenSet.Values() {
			token := wsToken.(WsToken)
			if time.Now().Unix()-token.CreateTime.Unix() > tokenEffectiveDuration {
				wsTokenSet.Remove(token)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func wsConnMonitor() {
	for {
		for i, e := range wsList.Values() {
			c := e.(*ws.Client)
			if time.Now().Unix()-c.LastUpdateTime.Unix() > wsTimeout {
				c.Close()
				wsList.Remove(i)
			}
		}

		log.Debugf("current wsList len is: %d", wsList.Size())
		time.Sleep(time.Duration(wsTimeout) * time.Second)
	}
}

func echoHandler() gin.HandlerFunc {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return func(c *gin.Context) {
		token := c.Query("token")

		tokenSet := hashset.New()
		for _, wsToken := range wsTokenSet.Values() {
			tokenSet.Add(wsToken.(WsToken).Token)
		}

		if !tokenSet.Contains(token) {
			c.JSON(419, result.Fail("login needed", ""))
			c.Abort()
			return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}
		portStructs := QueryPortStructsWithUsage()
		sort.Sort(sort.Reverse(PortStructSorter(portStructs)))

		conn.WriteJSON(portStructs)
		client := ws.Client{Conn: conn}
		wsList.Add(&client)

		go func() {
			for {
				if client.Closed() {
					break
				}
				_, _, err := client.Conn.ReadMessage()
				if err != nil {
					break
				}
				client.UpdateTime()
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

var whiteList = arraylist.New("/api/login", "/api/echo")

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if whiteList.Contains(c.Request.URL.Path) {
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
		p := PortStruct{}
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
		p := PortStruct{}
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
