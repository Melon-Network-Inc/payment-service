package processor

import "github.com/gin-gonic/gin"

const (
	ContextUserIdKey   = "UserID"
	ContextUsernameKey = "Username"
	ContextRoleKey     = "Role"
)

func GetUserID(c *gin.Context) string {
	return c.Request.Header.Get(ContextUserIdKey)
}

func GetUsername(c *gin.Context) string {
	return c.Request.Header.Get(ContextUsernameKey)
}

func GetRole(c *gin.Context) string {
	return c.Request.Header.Get(ContextRoleKey)
}
