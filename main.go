package main

import (
	"github.com/Bpazy/ssManager/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ok = http.StatusOK
const failed = http.StatusInternalServerError

func main() {
	engine := gin.Default()
	engine.Use(user.DefaultAuth())
	engine.Any(user.DefaultLoginUrl, user.Login())

	engine.Any("/", defaultHandler)
	engine.Any("/findSsConfig", findConfigHandler)
	engine.Any("/addPortPassword", addPortPasswordHandler)
	engine.Run(config.Addr)
}

func defaultHandler(c *gin.Context) {
	c.JSON(ok, "OK")
}

func addPortPasswordHandler(c *gin.Context) {
	ssConfig, err := findSsConfig()
	if err != nil {
		c.JSON(failed, err)
		return
	}
	port, _ := c.GetQuery("port")
	password, _ := c.GetQuery("password")
	ssConfig.PortPassword[port] = password
	err = saveSsConfig(ssConfig)
	if err != nil {
		c.JSON(failed, err)
		return
	}
	c.JSON(ok, "add success.")
}

func findConfigHandler(c *gin.Context) {
	ssConfig, err := findSsConfig()
	if err != nil {
		c.JSON(failed, err)
		return
	}
	c.JSON(ok, ssConfig)
}
