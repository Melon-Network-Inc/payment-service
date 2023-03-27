package taskq

import (
	"fmt"
	"github.com/Melon-Network-Inc/common/pkg/config"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"
)

type QueueManager interface {
	// RegisterTxnStatusQueue registers a queue for the txn status worker
	RegisterTxnStatusQueue(serverConfig config.ServiceConfig)
	// RegisterTxnStatusTask registers a task for the txn status worker
	RegisterTxnStatusTask(
		*gin.Context,
		entity.Transaction,
		func(ctx *gin.Context, txn entity.Transaction) error) error
	// Range iterates over all registered queues.
	Range(func(taskq.Queue) bool)
	// StartConsumers starts all registered queues.
	StartConsumers(ctx *gin.Context) error
	// StopConsumers stops all registered queues.
	StopConsumers() error
	// Close closes all registered queues.
	Close() error
}

type queueManager struct {
	factory taskq.Factory
	queues  map[string]taskq.Queue
}

// RegisterTxnStatusQueue registers a queue for the txn status worker
func (q queueManager) RegisterTxnStatusQueue(serverConfig config.ServiceConfig) {
	q.queues["txn-status-worker"] = q.factory.RegisterQueue(&taskq.QueueOptions{
		Name: "txn-status-worker",
		Redis: redis.NewClient(&redis.Options{
			Addr: serverConfig.CacheUrl,
		}),
	})
}

// RegisterTxnStatusTask registers a task for the txn status worker
func (q queueManager) RegisterTxnStatusTask(
	ctx *gin.Context,
	txn entity.Transaction,
	checkStatus func(ctx *gin.Context, txn entity.Transaction) error) error {
	task := taskq.RegisterTask(&taskq.TaskOptions{
		Name:    fmt.Sprintf("check-status-%s-%s", txn.Blockchain, txn.TxId),
		Handler: checkStatus,
	})
	message := task.WithArgs(ctx, txn)
	if err := q.queues["txn-status-worker"].Add(message); err != nil {
		return err
	}
	return nil
}

// Range iterates over all registered queues.
func (q queueManager) Range(fn func(taskq.Queue) bool) {
	q.factory.Range(fn)
}

// StartConsumers starts all registered queues.
func (q queueManager) StartConsumers(ctx *gin.Context) error {
	return q.factory.StartConsumers(ctx)
}

// StopConsumers stops all registered queues.
func (q queueManager) StopConsumers() error {
	return q.factory.StopConsumers()
}

// Close closes all registered queues.
func (q queueManager) Close() error {
	return q.factory.Close()
}

func NewTaskQueueManager(serverConfig config.ServiceConfig) QueueManager {
	qm := queueManager{
		factory: redisq.NewFactory(),
		queues:  make(map[string]taskq.Queue),
	}
	qm.RegisterTxnStatusQueue(serverConfig)
	return &qm
}

func NewQueueManager() QueueManager {
	return &queueManager{
		factory: redisq.NewFactory(),
		queues:  make(map[string]taskq.Queue),
	}
}
