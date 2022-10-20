package activity

import (
	"net/http"

	"github.com/Melon-Network-Inc/common/pkg/mwerrors"
	"github.com/Melon-Network-Inc/common/pkg/pagination"

	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.RouterGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	routes := r.Group("/activity")
	routes.GET("/query", res.QueryActivities)
	routes.GET("/", res.ListActivities)
}

type resource struct {
	service Service
	logger  log.Logger
}

// QueryActivities godoc
// @Summary      Query activities of an account
// @Description  Query activities of an account
// @ID           query-activities
// @Tags         activities
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param page query string false "page number"
// @Param per_page query string false "page size"
// @Accept       json
// @Produce      json
// @Success      200 {array} api.Post
// @Failure      400
// @Failure      401
// @Failure      404
// @Router       /activity/query [get]
func (r resource) QueryActivities(c *gin.Context) {
	ownerID, friendIDs, count, err := r.service.Count(c)
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	pages := pagination.NewFromRequest(c.Request, count)
	posts, err := r.service.Query(c, pages.Offset(), pages.Limit(), ownerID, friendIDs)
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	pages.Items = posts
	c.JSON(http.StatusOK, &pages)
}

// ListActivities godoc
// @Summary      List activities of an account
// @Description  List activities of an account
// @ID           list-activities
// @Tags         activities
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Accept       json
// @Produce      json
// @Success      200 {array} api.ActivityResponse
// @Failure      400
// @Failure      401
// @Failure      404
// @Router       /activity [get]
func (r resource) ListActivities(c *gin.Context) {
	activities, err := r.service.List(c)
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	c.JSON(http.StatusOK, &activities)
}
