package transaction

import (
	"errors"

	"github.com/Melon-Network-Inc/entity-repo/pkg/api"
	"github.com/Melon-Network-Inc/entity-repo/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	"github.com/Melon-Network-Inc/payment-service/pkg/processor"
	"github.com/Melon-Network-Inc/payment-service/pkg/utils"
	"github.com/gin-gonic/gin"
)

const NotAllowOperation = "cannot delete transaction record from the user who is not related to this transaction"

// Service encapsulates usecase logic for transactions.
type Service interface {
	Add(ctx *gin.Context, input api.AddTransactionRequest) (Transaction, error)
	Get(c *gin.Context, ID string) (Transaction, error)
	List(ctx *gin.Context) ([]Transaction, error)
	ListByUser(ctx *gin.Context, ID string, showPrivate bool) ([]Transaction, error)
	Update(ctx *gin.Context, ID string, input api.UpdateTransactionRequest) (Transaction, error)
	Delete(ctx *gin.Context, ID string) (Transaction, error)
}

// transaction represents the data about a transaction.
type Transaction struct {
	entity.Transaction
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new transaction service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Create creates a new transaction.
func (s service) Add(ctx *gin.Context, req api.AddTransactionRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	ownerID, err := utils.Int(processor.GetUserID(ctx))
	if err != nil {
		return Transaction{}, err
	}
	if req.SenderId != ownerID && req.ReceiverId != ownerID {
		return Transaction{}, errors.New(NotAllowOperation)
	}

	transaction, err := s.repo.Add(ctx, entity.Transaction{
		Name:           req.Name,
		Status:         req.Status,
		Amount:         req.Amount,
		Currency:       req.Currency,
		SenderId:       req.SenderId,
		SenderPubkey:   req.SenderPubkey,
		ReceiverId:     req.ReceiverId,
		ReceiverPubkey: req.ReceiverPubkey,
		IsPublic:       req.IsPublic,
		Message:        req.Message,
	})
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Get returns the transaction with the specified the transaction ID.
func (s service) Get(ctx *gin.Context, ID string) (Transaction, error) {
	uid, err := utils.Int(ID)
	if err != nil {
		return Transaction{}, err
	}
	transaction, err := s.repo.Get(ctx, uid)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Get returns the a list of transactions associated to the requester.
func (s service) List(ctx *gin.Context) ([]Transaction, error) {
	return s.ListByUser(ctx, processor.GetUserID(ctx), true)
}

// Get returns the a list of transactions associated to a user.
func (s service) ListByUser(ctx *gin.Context, ID string, showPrivate bool) ([]Transaction, error) {
	userID, err := utils.Int(ID)
	if err != nil {
		return []Transaction{}, err
	}

	transaction, err := s.repo.List(ctx, userID, showPrivate)
	if err != nil {
		return []Transaction{}, err
	}
	listTransaction := []Transaction{}
	for _, transaction := range transaction {
		listTransaction = append(listTransaction, Transaction{transaction})
	}
	return listTransaction, nil
}

// Update updates the transaction with the specified the transaction ID.
func (s service) Update(
	ctx *gin.Context,
	ID string,
	input api.UpdateTransactionRequest,
) (Transaction, error) {
	if err := input.Validate(); err != nil {
		return Transaction{}, err
	}

	UID, err := utils.Int(ID)
	if err != nil {
		return Transaction{}, err
	}

	transaction, err := s.repo.Get(ctx, UID)
	if err != nil {
		return Transaction{}, err
	}
	ownerID, err := utils.Uint(processor.GetUserID(ctx))
	if err != nil {
		return Transaction{}, err
	}

	if checkAllowsOperation(transaction, ownerID) {
		return Transaction{}, errors.New(NotAllowOperation)
	}

	if input.Name != "" {
		transaction.Name = input.Name
	}
	if input.Message != "" {
		transaction.Status = input.Message
	}
	if input.Status != "" {
		transaction.Status = input.Status
	}
	transaction.IsPublic = input.IsPublic

	if err := s.repo.Update(ctx, transaction); err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Delete deletes the transaction with the specified ID.
func (s service) Delete(ctx *gin.Context, ID string) (Transaction, error) {
	uid, err := utils.Int(ID)
	if err != nil {
		return Transaction{}, err
	}

	transaction, err := s.repo.Get(ctx, uid)
	if err != nil {
		return Transaction{}, err
	}
	ownerID, err := utils.Uint(processor.GetUserID(ctx))
	if err != nil {
		return Transaction{}, err
	}

	if checkAllowsOperation(transaction, ownerID) {
		return Transaction{}, errors.New(NotAllowOperation)
	}

	err = s.repo.Delete(ctx, transaction)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

func checkAllowsOperation(transaction entity.Transaction, ownerID uint) bool {
	return transaction.SenderId != int(ownerID) && transaction.ReceiverId != int(ownerID)
}
