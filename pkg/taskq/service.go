package taskq

import "github.com/Melon-Network-Inc/common/pkg/log"

// Service encapsulates use case logic for activities.
type Service interface {
	// Shutdown stops all registered queues.
	Shutdown() error
}

type service struct {
	queueManager QueueManager
	logger 	 	 log.Logger
}

// Shutdown stops all registered queues.
func (s service) Shutdown() error {
	return s.queueManager.StopConsumers()
}

// NewService creates a new service.
func NewService(queueManager QueueManager, logger log.Logger) Service {
	return service{
		queueManager: queueManager,
		logger: 	  logger,
	}
}
