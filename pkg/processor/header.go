package processor

import "github.com/gin-gonic/gin"

const (
	ContextUserIdKey = "User-Id"
	ContextUsernamekey = "Username"
)

func GetUserID(c *gin.Context) (string) {
   return c.Request.Header.Get(ContextUserIdKey)
}

func GetUsername(c *gin.Context) (string) {
   return c.Request.Header.Get(ContextUsernamekey)
}