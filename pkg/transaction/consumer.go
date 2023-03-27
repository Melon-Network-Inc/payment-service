package transaction

import (
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/gin-gonic/gin"
)

type Consumer interface {
	// CheckPendingTxns checks the status of the transaction.
	CheckPendingTxns(ctx *gin.Context, txn entity.Transaction) error
}


type consumer struct {
	service Service
	logger  log.Logger
}

// NewConsumer creates a new consumer.
func NewConsumer(service Service, logger log.Logger) Consumer {
	return consumer{
		service: service,
		logger:  logger,
	}
}

// CheckPendingTxns checks the status of the transaction.
func (c consumer) CheckPendingTxns(ctx *gin.Context, txn entity.Transaction) error {
	go func() {
		c.logger.Info("Start to check the status of the transaction", "txId", txn.TxId)
			
		taskQueue := *c.service.GetTaskQueueManager()
		// Check the status of the transaction.
		err := taskQueue.StartConsumers(ctx)
		if err != nil {
			return
		}
	}()
	return nil
}