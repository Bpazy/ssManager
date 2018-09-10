package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
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

func BindJson(c *gin.Context, target interface{}) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.Error(errors.New("json deserialize failed"))
		return false
	}
	return true
}
