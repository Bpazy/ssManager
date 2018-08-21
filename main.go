package main

import (
	"github.com/Bpazy/ssManager/user"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

const statusOk = http.StatusOK
const statusFailed = http.StatusInternalServerError

func main() {
	engine := gin.Default()

	engine.Use(user.DefaultAuth())
	engine.Any(user.DefaultLoginUrl, user.Login())

	engine.Any("/", defaultHandler)

	engine.Any("/findSsConfig", findConfigHandler)
	engine.Any("/addPortPassword", addPortPasswordHandler)
	engine.Any("/usage", usage)
	engine.Run(config.Addr)
}

func usage(c *gin.Context) {
	port, ok := c.GetQuery("port")
	if !ok {
		c.JSON(statusFailed, nil)
		return
	}
	command := "iptables -nvx -L | grep spt:{} | awk '{print $2}'"
	cmd := exec.Command("bash", "-c", strings.Replace(command, "{}", port, -1))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.JSON(statusFailed, err)
		return
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		c.JSON(statusFailed, err)
		return
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		c.JSON(statusFailed, err)
		return
	}
	c.JSON(statusOk, strings.TrimSpace(string(opBytes)))
}

func defaultHandler(c *gin.Context) {
	c.JSON(statusOk, "")
}

func addPortPasswordHandler(c *gin.Context) {
	ssConfig, err := findSsConfig()
	if err != nil {
		c.JSON(statusFailed, err)
		return
	}
	port, _ := c.GetQuery("port")
	password, _ := c.GetQuery("password")
	ssConfig.PortPassword[port] = password
	err = saveSsConfig(ssConfig)
	if err != nil {
		c.JSON(statusFailed, err)
		return
	}
	c.JSON(statusOk, "add success.")
}

func findConfigHandler(c *gin.Context) {
	ssConfig, err := findSsConfig()
	if err != nil {
		c.JSON(statusFailed, err)
		return
	}
	c.JSON(statusOk, ssConfig)
}
