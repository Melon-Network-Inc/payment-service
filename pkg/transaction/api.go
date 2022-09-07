package transaction

import (
	"net/http"

	"github.com/Melon-Network-Inc/common/pkg/api"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.RouterGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	routes := r.Group("/transactions")
	routes.POST("/", res.AddTransaction)
	routes.GET("/user/:id", res.GetAllTransactionsByUser)
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
// @Failure      404
// @Router       /transaction [post]
func (r resource) AddTransaction(c *gin.Context) {
	var input api.AddTransactionRequest
	// getting request's body
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	transaction, err := r.service.Add(c, input)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
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
// @Failure      404
// @Router       /transaction/{id} [get]
func (r resource) GetTransaction(c *gin.Context) {
	transaction, err := r.service.Get(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
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
// @Failure      404
// @Router       /transaction [get]
func (r resource) GetAllTransactions(c *gin.Context) {
	transactions, err := r.service.List(c)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, &transactions)
}

// GetAllTransactionsByUser    godoc
// @Summary      List all transactions by an account
// @Description  List all transactions by an account
// @ID           list-transactions-by-user
// @Tags         transactions
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param id path int true "Transaction ID"
// @Accept       json
// @Produce      json
// @Success      200 {array} api.TransactionResponse
// @Failure      404
// @Router       /transaction [get]
func (r resource) GetAllTransactionsByUser(c *gin.Context) {
	transactions, err := r.service.ListByUser(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
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
// @Failure      404
// @Router       /transaction [put]
func (r resource) UpdateTransaction(c *gin.Context) {
	var input api.UpdateTransactionRequest
	// getting request's body
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	transaction, err := r.service.Update(c, c.Param("id"), input)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
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
// @Failure      404
// @Router       /transaction [delete]
func (r resource) DeleteTransaction(c *gin.Context) {
	transaction, err := r.service.Delete(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, &transaction)
}
