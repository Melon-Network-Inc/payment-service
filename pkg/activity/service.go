package activity

import (
	"github.com/Melon-Network-Inc/payment-service/pkg/utils"
	"github.com/emirpasic/gods/sets/hashset"
	"sort"

	"github.com/Melon-Network-Inc/common/pkg/mwerrors"

	"github.com/Melon-Network-Inc/account-service/pkg/processor"
	accountRepo "github.com/Melon-Network-Inc/account-service/pkg/repository"
	"github.com/Melon-Network-Inc/common/pkg/api"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/Melon-Network-Inc/payment-service/pkg/repository"
	"github.com/gin-gonic/gin"
)

// Service encapsulates use case logic for activities.
type Service interface {
	List(c *gin.Context) (api.ActivityResponse, error)
	Count(c *gin.Context) (uint, []uint, int, error)
	Query(c *gin.Context, offset, limit int, ownerID uint, friendIDs []uint) ([]api.Post, error)
}

type service struct {
	userRepo        accountRepo.UserRepository
	transactionRepo repository.TransactionRepository
	friendRepo      accountRepo.FriendRepository
	logger          log.Logger
}

// NewService creates a new address service.
func NewService(
	userRepo accountRepo.UserRepository,
	transactionRepo repository.TransactionRepository,
	friendRepo accountRepo.FriendRepository,
	logger log.Logger) Service {
	return service{userRepo, transactionRepo, friendRepo, logger}
}

// Count returns all friend's activities count.
func (s service) Count(c *gin.Context) (uint, []uint, int, error) {
	ownerID, err := processor.GetContextUserID(c)
	if err != nil {
		return 0, []uint{}, 0, err
	}

	user, err := s.userRepo.Get(c, ownerID)
	if err != nil {
		return 0, []uint{}, 0, mwerrors.NewResourceNotFoundWithID(ownerID)
	}
	relations, err := s.friendRepo.ListAllRelations(c, user)
	if err != nil {
		return 0, []uint{}, 0, mwerrors.NewServerError(err)
	}
	var friendIDs []uint
	for _, relation := range relations {
		friendIDs = append(friendIDs, relation.ToUserRef)
	}

	cnt, err := s.transactionRepo.CountByFriendIDs(c, ownerID, friendIDs)
	if err != nil {
		return 0, []uint{}, 0, mwerrors.NewServerError(err)
	}
	return ownerID, friendIDs, cnt, nil
}

// Query returns all friend's activities by page.
func (s service) Query(c *gin.Context,
	offset, limit int,
	ownerID uint,
	friendIDs []uint) ([]api.Post, error) {
	items, err := s.transactionRepo.QueryByFriendIDs(c, offset, limit, ownerID, friendIDs)
	if err != nil {
		return []api.Post{}, mwerrors.NewResourcesNotFound(err)
	}

	convertedTxns, err := s.ConvertToApiTransactions(c, items, true)
	if err != nil {
		return []api.Post{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	var posts []api.Post
	for _, txn := range convertedTxns {
		posts = append(posts, api.Post{
			Type:        api.TransactionPostType,
			Transaction: txn.Transaction,
			Moment:      api.Moment{},
		})
	}
	return posts, nil
}

// List returns all friend's activities with the specified the user ID.
func (s service) List(c *gin.Context) (api.ActivityResponse, error) {
	ownerID, err := processor.GetContextUserID(c)
	if err != nil {
		return api.ActivityResponse{}, err
	}

	user, err := s.userRepo.Get(c, ownerID)
	if err != nil {
		return api.ActivityResponse{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}
	relations, err := s.friendRepo.ListAllRelations(c, user)
	if err != nil {
		return api.ActivityResponse{}, mwerrors.NewResourcesNotFound(err)
	}

	// Query all friends' activities
	var txnsActivities []entity.Transaction
	for _, relation := range relations {
		transactions, err := s.transactionRepo.ListByUserID(c, user.ID, relation.ToUserRef)
		if err != nil {
			return api.ActivityResponse{}, mwerrors.NewResourcesNotFound(err)
		}
		txnsActivities = append(txnsActivities, transactions...)
	}

	// Query requester activities
	transactions, err := s.transactionRepo.ListByRequester(c, ownerID)
	if err != nil {
		return api.ActivityResponse{}, mwerrors.NewResourcesNotFound(err)
	}
	txnsActivities = append(txnsActivities, transactions...)

	// Sort transaction activities by updated at.
	sort.Slice(txnsActivities, func(i, j int) bool {
		return txnsActivities[i].UpdatedAt.Before(txnsActivities[j].UpdatedAt)
	})

	// Convert entity transaction to api transactions.
	convertedTxns, err := s.ConvertToApiTransactions(c, txnsActivities, true)
	if err != nil {
		return api.ActivityResponse{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	// Add transaction to Post
	var posts []api.Post
	for _, txn := range convertedTxns {
		posts = append(posts, api.Post{
			Type:        api.TransactionPostType,
			Transaction: txn.Transaction,
			Moment:      api.Moment{},
		})
	}
	return api.ActivityResponse{Posts: posts}, nil
}

func (s service) ConvertToApiTransactions(c *gin.Context, txns []entity.Transaction, isPrune bool) ([]api.TransactionResponse, error) {
	userMap := make(map[uint]entity.User)
	userIDSet := hashset.New()

	for _, txn := range txns {
		userIDSet.Add(txn.SenderId)
		userIDSet.Add(txn.ReceiverId)
	}

	users, exists, err := s.userRepo.GetByIDs(c, utils.GetUints(userIDSet.Values()))
	if err != nil {
		return []api.TransactionResponse{}, err
	}
	if !exists {
		return []api.TransactionResponse{}, nil
	}
	for _, user := range users {
		userMap[user.ID] = user
	}

	var result []api.TransactionResponse
	for _, txn := range txns {
		sender := userMap[uint(txn.SenderId)]
		receiver := userMap[uint(txn.ReceiverId)]
		result = append(result, convert(txn, sender, receiver, isPrune))
	}
	return result, nil
}

func convert(txn entity.Transaction, sender, receiver entity.User, prune bool) api.TransactionResponse {
	convertedTxn := api.Transaction{
		ID:               int(txn.ID),
		Name:             txn.Name,
		Status:           txn.Status,
		Amount:           "",
		Currency:         txn.Currency,
		Blockchain:       txn.Blockchain,
		Symbol:           txn.Symbol,
		SenderID:         txn.SenderId,
		SenderUsername:   sender.Username,
		SenderUrl:        sender.Avatar,
		SenderPubkey:     "",
		ReceiverID:       txn.ReceiverId,
		ReceiverUsername: receiver.Username,
		ReceiverUrl:      receiver.Avatar,
		ReceiverPubkey:   "",
		TransactionType:  txn.TransactionType,
		Message:          txn.Message,
	}
	if !prune {
		convertedTxn.Amount = utils.GetFloatPointString(txn.Amount)
		convertedTxn.SenderPubkey = txn.SenderPubkey
		convertedTxn.ReceiverPubkey = txn.ReceiverPubkey
	}
	return api.TransactionResponse{Transaction: convertedTxn}
}
