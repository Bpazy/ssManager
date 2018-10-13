package util

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
)

func ShouldParseInt64(i string) (int64, bool) {
	atoi, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		return 0, false
	}
	return atoi, true
}

func MustParseInt64(i string) int64 {
	atoi, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		panic(err)
	}
	return atoi
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
	if err != nil {
		panic(err)
	}
	return result
}

func RunCommand(c string) (string, error) {
	log.Debugf("command prepare: " + c)
	cmd := exec.Command("bash", "-c", c)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	result := strings.TrimSpace(out.String())
	log.Debugf("command result: " + result)
	return result, err
}
