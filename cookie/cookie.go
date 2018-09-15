package cookie

import (
	"github.com/Bpazy/ssManager/aes"
	"github.com/Bpazy/ssManager/base64"
	"github.com/gin-gonic/gin"
)

const (
	cookieName = "cookieName"
)

var (
	key = []byte("1234567890123456")
)

func SaveUserId(c *gin.Context, userId string) {
	encryptedUserId, err := aes.Encrypt([]byte(userId), key)
	if err != nil {
		panic(err)
	}
	c.SetCookie(cookieName, base64.Encode(encryptedUserId), 60*60*8, "/", "", false, false)
}

func GetSavedUserId(c *gin.Context) (string, error) {
	base64UserId, err := c.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	encryptedUserId, err := base64.Decode(base64UserId)
	if err != nil {
		return "", err
	}
	userId, err := aes.Decrypt(encryptedUserId, key)
	if err != nil {
		return "", nil
	}
	return string(userId), nil
}
