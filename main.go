package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
)

type SsConfig struct {
	Server       string            `json:"server"`
	LocalAddress string            `json:"local_address"`
	LocalPort    int               `json:"local_port"`
	Timeout      int               `json:"timeout"`
	Method       string            `json:"method"`
	FastOpen     string            `json:"fast_open"`
	PortPassword map[string]string `json:"port_password"`
}

const ok = 200
const failed = 500

func main() {
	r := gin.Default()
	r.GET("/findConfig", func(c *gin.Context) {
		ssConfig, err := findConfig()
		if err != nil {
			c.JSON(failed, err)
			return
		}
		c.JSON(ok, ssConfig)
	})
	r.Run()
}

func findConfig() (*SsConfig, error) {
	fileBytes, err := ioutil.ReadFile("./test.json")
	if err != nil {
		return nil, err
	}
	ssConfig := SsConfig{}
	err = json.Unmarshal(fileBytes, &ssConfig)
	if err != nil {
		return nil, err
	}
	return &ssConfig, nil
}

func test() {
	fileBytes, err := ioutil.ReadFile("./test.json")
	if err != nil {
		log.Fatal(err)
	}
	ssConfig := SsConfig{}
	err = json.Unmarshal(fileBytes, &ssConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ssConfig)

	ssConfig.PortPassword["9999"] = "hahaha"

	ssConfigBytes, err := json.Marshal(ssConfig)
	if err != nil {
		log.Fatal(err)
	}

	ssConfigBytes, err = prettyPrint(ssConfigBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./test2.json", ssConfigBytes, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
