package activity

import (
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
	var posts []api.Post
	for _, item := range items {
		posts = append(posts, api.Post{
			Type:        api.TransactionPostType,
			Transaction: item,
			Moment:      entity.Moment{},
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
	var transactionsHistory []entity.Transaction
	for _, relation := range relations {
		transactions, err := s.transactionRepo.ListByUserID(c, user.ID, relation.ToUserRef)
		if err != nil {
			return api.ActivityResponse{}, mwerrors.NewResourcesNotFound(err)
		}
		transactionsHistory = append(transactionsHistory, transactions...)
	}

	// Query requester activities
	transactions, err := s.transactionRepo.ListByRequester(c, ownerID)
	if err != nil {
		return api.ActivityResponse{}, mwerrors.NewResourcesNotFound(err)
	}
	transactionsHistory = append(transactionsHistory, transactions...)

	sort.Slice(transactionsHistory, func(i, j int) bool {
		return transactionsHistory[i].UpdatedAt.Before(transactionsHistory[j].UpdatedAt)
	})

	var posts []api.Post
	for _, transaction := range transactionsHistory {
		posts = append(posts, api.Post{
			Type:        api.TransactionPostType,
			Transaction: transaction,
			Moment:      entity.Moment{},
		})
	}
	return api.ActivityResponse{Posts: posts}, nil
}
