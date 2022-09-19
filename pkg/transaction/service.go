package transaction

import (
	"errors"

	accountRepo "github.com/Melon-Network-Inc/account-service/pkg/repository"
	"github.com/Melon-Network-Inc/payment-service/pkg/repository"

	"github.com/Melon-Network-Inc/common/pkg/api"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"

	"github.com/Melon-Network-Inc/payment-service/pkg/processor"
	"github.com/Melon-Network-Inc/payment-service/pkg/utils"
	"github.com/gin-gonic/gin"
)

const NotAllowOperation = "requester does not have right access to change target resource"

// Service encapsulates use case logic for transactions.
type Service interface {
	Add(ctx *gin.Context, input api.AddTransactionRequest) (api.TransactionResponse, error)
	Get(c *gin.Context, ID string) (api.TransactionResponse, error)
	List(ctx *gin.Context) ([]api.TransactionResponse, error)
	ListByUser(ctx *gin.Context, ID string) ([]api.TransactionResponse, error)
	ListByUserWithShowType(ctx *gin.Context, ID string, showType string) ([]api.TransactionResponse, error)
	Update(ctx *gin.Context, ID string, input api.UpdateTransactionRequest) (api.TransactionResponse, error)
	Delete(ctx *gin.Context, ID string) (api.TransactionResponse, error)
}

type service struct {
	transactionRepo repository.TransactionRepository
	userRepo        accountRepo.UserRepository
	friendRepo      accountRepo.FriendRepository
	logger          log.Logger
}

// NewService creates a new transaction service.
func NewService(
	transactionRepo repository.TransactionRepository,
	userRepo accountRepo.UserRepository,
	friendRepo accountRepo.FriendRepository,
	logger log.Logger) Service {
	return service{transactionRepo, userRepo, friendRepo, logger}
}

// Add creates a new transaction.
func (s service) Add(ctx *gin.Context, req api.AddTransactionRequest) (api.TransactionResponse, error) {
	if err := req.Validate(); err != nil {
		return api.TransactionResponse{}, err
	}

	ownerID, err := utils.Int(processor.GetUserID(ctx))
	if err != nil {
		return api.TransactionResponse{}, err
	}
	if req.SenderId != ownerID && req.ReceiverId != ownerID {
		return api.TransactionResponse{}, errors.New(NotAllowOperation)
	}

	transaction, err := s.transactionRepo.Add(ctx, entity.Transaction{
		Name:           req.Name,
		Status:         req.Status,
		Amount:         req.Amount,
		Currency:       req.Currency,
		SenderId:       req.SenderId,
		SenderPubkey:   req.SenderPubkey,
		ReceiverId:     req.ReceiverId,
		ReceiverPubkey: req.ReceiverPubkey,
		ShowType:       req.ShowType,
		Message:        req.Message,
	})
	if err != nil {
		return api.TransactionResponse{}, err
	}
	return api.TransactionResponse{Transaction: transaction}, nil
}

// Get returns the transaction with the specified the transaction ID.
func (s service) Get(ctx *gin.Context, ID string) (api.TransactionResponse, error) {
	uid, err := utils.Int(ID)
	if err != nil {
		return api.TransactionResponse{}, err
	}
	transaction, err := s.transactionRepo.Get(ctx, uid)
	if err != nil {
		return api.TransactionResponse{}, err
	}
	return api.TransactionResponse{Transaction: transaction}, nil
}

// List returns the a list of transactions associated to the requester.
func (s service) List(ctx *gin.Context) ([]api.TransactionResponse, error) {
	return s.ListByUserWithShowType(ctx, processor.GetUserID(ctx), "Private")
}

// ListByUser returns the a list of transactions associated to target user depending on requester's relation.
func (s service) ListByUser(ctx *gin.Context, ID string) ([]api.TransactionResponse, error) {
	userID := processor.GetUserID(ctx)
	if userID == "" {
		return []api.TransactionResponse{}, errors.New("missing request user information")
	}
	if userID == ID {
		return s.List(ctx)
	}

	requesterID, err := utils.Uint(userID)
	if err != nil {
		return []api.TransactionResponse{}, err
	}
	requestUser, err := s.userRepo.Get(ctx, requesterID)
	if err != nil {
		return []api.TransactionResponse{}, err
	}

	otherID, err := utils.Uint(ID)
	if err != nil {
		return []api.TransactionResponse{}, err
	}
	otherUser, err := s.userRepo.Get(ctx, otherID)
	if err != nil {
		return []api.TransactionResponse{}, err
	}

	showType := "Public"
	exists, err := s.friendRepo.HasRelationByBothUsers(ctx, requestUser, otherUser)
	if err != nil {
		return []api.TransactionResponse{}, err
	}
	if exists {
		showType = "Friend"
	}

	return s.ListByUserWithShowType(ctx, ID, showType)
}

// ListByUserWithShowType returns the a list of transactions associated to a user.
func (s service) ListByUserWithShowType(ctx *gin.Context, ID string, showType string) ([]api.TransactionResponse, error) {
	userID, err := utils.Int(ID)
	if err != nil {
		return []api.TransactionResponse{}, err
	}

	transaction, err := s.transactionRepo.List(ctx, userID, showType)
	if err != nil {
		return []api.TransactionResponse{}, err
	}
	var listTransaction []api.TransactionResponse
	for _, transaction := range transaction {
		listTransaction = append(listTransaction, api.TransactionResponse{Transaction: transaction})
	}
	return listTransaction, nil
}

// Update updates the transaction with the specified the transaction ID.
func (s service) Update(
	ctx *gin.Context,
	ID string,
	input api.UpdateTransactionRequest,
) (api.TransactionResponse, error) {
	if err := input.Validate(); err != nil {
		return api.TransactionResponse{}, err
	}

	UID, err := utils.Int(ID)
	if err != nil {
		return api.TransactionResponse{}, err
	}

	transaction, err := s.transactionRepo.Get(ctx, UID)
	if err != nil {
		return api.TransactionResponse{}, err
	}
	ownerID, err := utils.Uint(processor.GetUserID(ctx))
	if err != nil {
		return api.TransactionResponse{}, err
	}

	if checkAllowsOperation(transaction, ownerID) {
		return api.TransactionResponse{}, errors.New(NotAllowOperation)
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
	if input.ShowType != "" {
		transaction.ShowType = input.ShowType
	}

	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return api.TransactionResponse{}, err
	}
	return api.TransactionResponse{Transaction: transaction}, nil
}

// Delete deletes the transaction with the specified ID.
func (s service) Delete(ctx *gin.Context, ID string) (api.TransactionResponse, error) {
	uid, err := utils.Int(ID)
	if err != nil {
		return api.TransactionResponse{}, err
	}

	transaction, err := s.transactionRepo.Get(ctx, uid)
	if err != nil {
		return api.TransactionResponse{}, err
	}
	ownerID, err := utils.Uint(processor.GetUserID(ctx))
	if err != nil {
		return api.TransactionResponse{}, err
	}

	if checkAllowsOperation(transaction, ownerID) {
		return api.TransactionResponse{}, errors.New(NotAllowOperation)
	}

	err = s.transactionRepo.Delete(ctx, transaction)
	if err != nil {
		return api.TransactionResponse{}, err
	}
	return api.TransactionResponse{Transaction: transaction}, nil
}

func checkAllowsOperation(transaction entity.Transaction, ownerID uint) bool {
	return transaction.SenderId != int(ownerID) && transaction.ReceiverId != int(ownerID)
}
