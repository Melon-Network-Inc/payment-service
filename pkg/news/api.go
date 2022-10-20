package news

import (
	"net/http"

	"github.com/Melon-Network-Inc/common/pkg/mwerrors"
	"github.com/Melon-Network-Inc/common/pkg/pagination"

	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.RouterGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	routes := r.Group("/news")
	routes.GET("/query", res.QueryActivities)
}

type resource struct {
	service Service
	logger  log.Logger
}

// QueryNews godoc
// @Summary      Query news by page
// @Description  Query news by page
// @ID           query-news
// @Tags         news
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param page query string false "page number"
// @Param per_page query string false "page size"
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.News
// @Failure      400
// @Failure      401
// @Failure      404
// @Router       /news/query [get]
func (r resource) QueryActivities(c *gin.Context) {
	count, err := r.service.Count(c)
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	pages := pagination.NewFromRequest(c.Request, count)
	news, err := r.service.Query(c, pages.Offset(), pages.Limit())
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	pages.Items = news
	c.JSON(http.StatusOK, &pages)
}