package taskq

import (
	"net/http"

	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/Melon-Network-Inc/common/pkg/mwerrors"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.RouterGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	routes := r.Group("/task")
	routes.GET("/shutdown", res.ShutdownTaskQueue)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) ShutdownTaskQueue(c *gin.Context) {
	r.logger.Debug("ShutdownTask")
	if err := r.service.Shutdown(); err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}