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

func MustParseInt(i string) int {
	atoi, e := strconv.Atoi(i)
	ShouldPanic(e)
	return atoi
}

func BindJson(c *gin.Context, target interface{}) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.Error(errors.New("json deserialize failed"))
		return false
	}
	return true
}
