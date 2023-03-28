package transaction

import (
	"context"
	"time"

	"github.com/Melon-Network-Inc/common/pkg/log"
)

type Consumer interface {
	// CheckPendingTxns checks the status of the pending transactions.
	CheckPendingTxns() error
}


type consumer struct {
	service Service
	logger  log.Logger
}

func NewConsumer(service Service, logger log.Logger) Consumer {
	return consumer{
		service: service,
		logger:  logger,
	}
}

// CheckPendingTxns checks the status of the transaction.
func (c consumer) CheckPendingTxns() error {
	// Wait for 30 seconds to check the status of the transaction.
	time.Sleep(2 * time.Second)
	c.logger.Info("Start to check the status of the pending transactions")
		
	taskQueue := *c.service.GetTaskQueueManager()
	// Check the status of the transaction.
	err := taskQueue.StartConsumers(context.Background())
	if err != nil {
		return err
	}
	return nil
}