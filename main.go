package main

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

const ok = http.StatusOK
const failed = http.StatusInternalServerError

func main() {
	//engine := gin.Default()
	//
	//engine.Use(user.DefaultAuth())
	//engine.Any(user.DefaultLoginUrl, user.Login())
	//
	//engine.Any("/", defaultHandler)
	//
	//engine.Any("/findSsConfig", findConfigHandler)
	//engine.Any("/addPortPassword", addPortPasswordHandler)
	//engine.Run(config.Addr)
	test()
}

func defaultHandler(c *gin.Context) {
	c.JSON(ok, "")
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

func test() {
	cmd := exec.Command("bash", "-c", "iptables -")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(opBytes))
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}
