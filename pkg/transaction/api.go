package transaction

import (
	"net/http"

	"github.com/Melon-Network-Inc/common/pkg/api"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/Melon-Network-Inc/common/pkg/mwerrors"
	"github.com/Melon-Network-Inc/common/pkg/pagination"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.RouterGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	routes := r.Group("/transaction")
	routes.POST("/", res.AddTransaction)
	routes.GET("/user/:id", res.GetAllTransactionsByUser)
	routes.GET("/query/:id", res.QueryTransactions)
	routes.GET("/:id", res.GetTransaction)
	routes.GET("/", res.GetAllTransactions)
	routes.PUT("/:id", res.UpdateTransaction)
	routes.DELETE("/:id", res.DeleteTransaction)
}

type resource struct {
	service Service
	logger  log.Logger
}

// AddTransaction    godoc
// @Summary      Add a transaction to account
// @Description  Add a transaction to account
// @ID           add-transaction
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param transaction body api.AddTransactionRequest true "Transaction Data"
// @Accept       json
// @Produce      json
// @Success      201 {object} api.TransactionResponse
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Router       /transaction [post]
func (r resource) AddTransaction(c *gin.Context) {
	var input api.AddTransactionRequest
	// getting request's body
	if err := c.BindJSON(&input); err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	r.logger.Debug("AddTransaction", input)
	transaction, err := r.service.Add(c, input)
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	c.JSON(http.StatusCreated, &transaction)
}

// GetTransaction    godoc
// @Summary      Get a transaction
// @Description  Get a transaction
// @ID           get-transaction
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param id path int true "Transaction ID"
// @Accept       json
// @Produce      json
// @Success      200 {object} api.TransactionResponse
// @Failure      400
// @Failure      401
// @Failure      404
// @Router       /transaction/{id} [get]
func (r resource) GetTransaction(c *gin.Context) {
	transaction, err := r.service.Get(c, c.Param("id"))
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	c.JSON(http.StatusOK, &transaction)
}

// GetAllTransactions    godoc
// @Summary      List all transactions of requester
// @Description  List all transactions of requester
// @ID           list-transactions
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Accept       json
// @Produce      json
// @Success      200 {array} api.TransactionResponse
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Router       /transaction [get]
func (r resource) GetAllTransactions(c *gin.Context) {
	transactions, err := r.service.List(c)
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	c.JSON(http.StatusOK, &transactions)
}

// QueryTransactions    godoc
// @Summary      Query requester's transactions by page
// @Description  Query requester's transactions by page
// @ID           query-transactions
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param id path int true "User ID"
// @Param page query string false "page number"
// @Param per_page query string false "page size"
// @Accept       json
// @Produce      json
// @Success      200 {array} api.TransactionResponse
// @Failure      400
// @Failure      401
// @Failure      404
// @Router       /transaction/query/{id} [get]
func (r resource) QueryTransactions(c *gin.Context) {
	showType, count, err := r.service.CountByUser(c, c.Param("id"))
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	pages := pagination.NewFromRequest(c.Request, count)
	addresses, err := r.service.Query(c, c.Param("id"), showType, pages.Offset(), pages.Limit())
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	pages.Items = addresses
	c.JSON(http.StatusOK, &pages)
}

// GetAllTransactionsByUser    godoc
// @Summary      List all transactions of an account
// @Description  List all transactions of an account
// @ID           list-transactions-by-user
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param id path int true "User ID"
// @Accept       json
// @Produce      json
// @Success      200 {array} api.TransactionResponse
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Router       /transaction/user/{id} [get]
func (r resource) GetAllTransactionsByUser(c *gin.Context) {
	transactions, err := r.service.ListByUser(c, c.Param("id"))
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	c.JSON(http.StatusOK, &transactions)
}

// UpdateTransaction godoc
// @Summary      Update a transaction
// @Description  Update a transaction
// @ID           update-transaction
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param transaction body api.UpdateTransactionRequest true "Transaction Data"
// @Accept       json
// @Produce      json
// @Success      200 {object} api.TransactionResponse
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      404
// @Failure      500
// @Router       /transaction [put]
func (r resource) UpdateTransaction(c *gin.Context) {
	var input api.UpdateTransactionRequest
	// getting request's body
	if err := c.BindJSON(&input); err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	transaction, err := r.service.Update(c, c.Param("id"), input)
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}

	c.JSON(http.StatusOK, &transaction)
}

// DeleteTransaction godoc
// @Summary      Delete a transaction
// @Description  Delete a transaction
// @ID           delete-transaction
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param id path int true "Transaction ID"
// @Accept       json
// @Produce      json
// @Success      200 {object} api.TransactionResponse
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Router       /transaction [delete]
func (r resource) DeleteTransaction(c *gin.Context) {
	transaction, err := r.service.Delete(c, c.Param("id"))
	if err != nil {
		mwerrors.HandleErrorResponse(c, r.logger, err)
		return
	}
	c.JSON(http.StatusOK, &transaction)
}
