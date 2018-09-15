package util

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func ShouldPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ShouldParseInt64(i string) (int64, bool) {
	atoi, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		return 0, false
	}
	return atoi, true
}

func MustParseInt(i string) int {
	atoi, e := strconv.Atoi(i)
	if e != nil {
		panic(e)
	}
	return atoi
}

func BindJson(c *gin.Context, target interface{}) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.Error(errors.New("json deserialize failed"))
		return false
	}
	return true
}

func MustRunCommand(c string) string {
	result, err := RunCommand(c)
	ShouldPanic(err)
	return result
}

func RunCommand(c string) (string, error) {
	log.Println("command prepare: " + c)
	cmd := exec.Command("bash", "-c", c)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	result := strings.TrimSpace(out.String())
	log.Println("command result: " + result)
	return result, err
}
