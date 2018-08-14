package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	userId   string
	userName string
	mobile   string
}

const DefaultLoginUrl = "/login"

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.GetQuery("userId")
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.SetCookie("userId", userId, 60*60, "/", ".", false, false)
		c.JSON(http.StatusOK, "Login success")
	}
}

func DefaultAuth() gin.HandlerFunc {
	return Auth(DefaultLoginUrl)
}

func Auth(loginUrl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == loginUrl {
			return
		}
		user, found := GetUser(c)
		if !found {
			c.JSON(http.StatusOK, "Login needed.")
			c.Abort()
			return
		}
		c.Set("user", user)
	}
}

func GetUser(c *gin.Context) (*User, bool) {
	userId, err := c.Cookie("userId")
	if err != nil {
		return nil, false
	}
	return &User{userId: userId, userName: "DefaultUserName", mobile: "default user mobile"}, true
}
