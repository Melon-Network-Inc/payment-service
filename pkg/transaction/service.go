package transaction

import (
	"fmt"
	"reflect"
	"strconv"

	accountRepo "github.com/Melon-Network-Inc/account-service/pkg/repository"
	"github.com/Melon-Network-Inc/payment-service/pkg/repository"
	"github.com/emirpasic/gods/sets/hashset"

	"github.com/Melon-Network-Inc/common/pkg/api"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/Melon-Network-Inc/common/pkg/mwerrors"
	"github.com/Melon-Network-Inc/common/pkg/notification"

	"github.com/Melon-Network-Inc/payment-service/pkg/processor"
	"github.com/Melon-Network-Inc/payment-service/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Service encapsulates use case logic for transactions.
type Service interface {
	Add(ctx *gin.Context, input api.AddTransactionRequest) (api.TransactionResponse, error)
	Get(c *gin.Context, ID string) (api.TransactionResponse, error)
	List(ctx *gin.Context) ([]api.TransactionResponse, error)
	ListByUser(ctx *gin.Context, ID string) ([]api.TransactionResponse, error)
	ListByUserWithShowType(ctx *gin.Context, ID string, showType string) ([]api.TransactionResponse, error)
	Update(ctx *gin.Context, ID string, input api.UpdateTransactionRequest) (api.TransactionResponse, error)
	Delete(ctx *gin.Context, ID string) (api.TransactionResponse, error)
	Count(c *gin.Context) (string, int, error)
	CountByUser(c *gin.Context, ID string) (string, int, error)
	CountByUserWithShowType(c *gin.Context, ID string, showType string) (string, int, error)
	Query(c *gin.Context, ID, showType string, offset, limit int) ([]api.TransactionResponse, error)
}

type service struct {
	transactionRepo  repository.TransactionRepository
	userRepo         accountRepo.UserRepository
	friendRepo       accountRepo.FriendRepository
	deviceRepo       accountRepo.DeviceRepository
	notificationRepo accountRepo.NotificationRepository
	fcmClient        *notification.FCMClient
	logger           log.Logger
}

// NewService creates a new transaction service.
func NewService(
	transactionRepo repository.TransactionRepository,
	userRepo accountRepo.UserRepository,
	friendRepo accountRepo.FriendRepository,
	deviceRepo accountRepo.DeviceRepository,
	notificationRepo accountRepo.NotificationRepository,
	fcmClient *notification.FCMClient,
	logger log.Logger) Service {
	return service{transactionRepo, userRepo, friendRepo, deviceRepo, notificationRepo, fcmClient, logger}
}

// Add creates a new transaction.
func (s service) Add(ctx *gin.Context, req api.AddTransactionRequest) (api.TransactionResponse, error) {
	if err := req.Validate(); err != nil {
		return api.TransactionResponse{}, mwerrors.NewIllegalArgumentError(err)
	}

	userID := processor.GetUserID(ctx)
	if userID == "" {
		return api.TransactionResponse{}, mwerrors.NewMissingAuthToken()
	}

	ownerID, err := utils.Int(userID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewIllegalArgumentError(err)
	}

	if req.SenderId != ownerID && req.ReceiverId != ownerID {
		return api.TransactionResponse{}, mwerrors.NewResourceNotAllowedWithOnlyUsername(processor.GetUsername(ctx))
	}

	txn := entity.Transaction{
		Name:           req.Name,
		Status:         req.Status,
		Amount:         req.Amount,
		Symbol:         req.Symbol,
		Blockchain:     req.Blockchain,
		SenderId:       req.SenderId,
		SenderPubkey:   req.SenderPubkey,
		ReceiverId:     req.ReceiverId,
		ReceiverPubkey: req.ReceiverPubkey,
		ShowType:       req.ShowType,
		Message:        req.Message,
	}
	if req.Currency != "" {
		txn.Currency = req.Currency
	}
	if req.TransactionType != "" {
		txn.TransactionType = req.TransactionType
	} else {
		txn.TransactionType = "standard"
	}
	createdTxn, err := s.transactionRepo.Add(ctx, txn)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewServerError(err)
	}

	user, err := s.userRepo.Get(ctx, uint(req.SenderId))
	if err != nil {
		return convert(createdTxn, entity.User{}, entity.User{}, false), mwerrors.NewResourcesNotFound(err)
	}
	otherUser, err := s.userRepo.Get(ctx, uint(req.ReceiverId))
	if err != nil {
		return convert(createdTxn, user, entity.User{}, false), mwerrors.NewResourcesNotFound(err)
	}
	devices, err := s.deviceRepo.GetDevices(ctx, otherUser)
	if err != nil {
		return convert(createdTxn, user, otherUser, false), mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	var aggregatedDevices string
	var tokenList []string
	if len(devices) != 0 {
		aggregatedDevices, tokenList = extractDeviceNameAndToken(devices)
	}

	newNotification := entity.Notification{
		UserRef:    uint(req.ReceiverId),
		ActorRef:   uint(req.SenderId),
		Device:     aggregatedDevices,
		Type:       entity.TransactionConfirmationType,
		Actor:      entity.ActorUserType,
		Title:      "Transaction Notification",
		Message:    CreateTransactionMessage(user, otherUser, createdTxn),
		TemplateID: 1,
	}
	createNotification, err := s.notificationRepo.CreateNotification(ctx, newNotification)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewServerError(err)
	}

	if len(devices) == 0 {
		return convert(createdTxn, user, otherUser, false), nil
	}
	if err = s.fcmClient.NotifyDevices(ctx, tokenList, createNotification); err != nil {
		return convert(createdTxn, user, otherUser, false), mwerrors.NewServerError(err)
	}

	return convert(createdTxn, user, otherUser, false), nil
}

// Get returns the transaction with the specified the transaction ID.
func (s service) Get(ctx *gin.Context, ID string) (api.TransactionResponse, error) {
	userID := processor.GetUserID(ctx)
	if userID == "" {
		return api.TransactionResponse{}, mwerrors.NewMissingAuthToken()
	}

	UID, err := utils.Uint(ID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewIllegalInputErrorWithMessage(err.Error())
	}

	transaction, err := s.transactionRepo.Get(ctx, UID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	txn, err := s.ConvertToApiTransaction(ctx, transaction, userID == ID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	return txn, nil
}

// List returns the list of transactions associated to the requester.
func (s service) List(ctx *gin.Context) ([]api.TransactionResponse, error) {
	return s.ListByUserWithShowType(ctx, processor.GetUserID(ctx), "Private")
}

// ListByUser returns the list of transactions associated to target user depending on requester's relation.
func (s service) ListByUser(ctx *gin.Context, ID string) ([]api.TransactionResponse, error) {
	userID := processor.GetUserID(ctx)
	if userID == "" {
		return []api.TransactionResponse{}, mwerrors.NewMissingAuthToken()
	}

	if userID == ID {
		return s.List(ctx)
	}

	requesterID, err := utils.Uint(userID)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewInvalidAuthToken(err)
	}
	requestUser, err := s.userRepo.Get(ctx, requesterID)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	otherID, err := utils.Uint(ID)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewIllegalArgumentError(err)
	}
	otherUser, err := s.userRepo.Get(ctx, otherID)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	showType := "Public"
	exists, err := s.friendRepo.HasRelationByBothUsers(ctx, requestUser, otherUser)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewResourceNotFoundWithPublicError(err)
	}
	if exists {
		showType = "Friend"
	}

	return s.ListByUserWithShowType(ctx, ID, showType)
}

// ListByUserWithShowType returns the list of transactions associated to a user.
func (s service) ListByUserWithShowType(ctx *gin.Context, ID string, showType string) ([]api.TransactionResponse, error) {
	userID, err := utils.Uint(ID)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewMissingAuthToken()
	}

	txns, err := s.transactionRepo.List(ctx, userID, showType)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewResourcesNotFound(err)
	}

	resp, err := s.ConvertToApiTransactions(ctx, txns, showType != "Private")
	if err != nil {
		return []api.TransactionResponse{}, err
	}
	return resp, nil
}

// Update updates the transaction with the specified the transaction ID.
func (s service) Update(
	ctx *gin.Context,
	ID string,
	input api.UpdateTransactionRequest,
) (api.TransactionResponse, error) {
	if err := input.Validate(); err != nil {
		return api.TransactionResponse{}, mwerrors.NewIllegalArgumentError(err)
	}
	UID, err := utils.Uint(ID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewIllegalInputErrorWithMessage(err.Error())
	}
	ownerID, err := utils.Uint(processor.GetUserID(ctx))
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewInvalidAuthToken(err)
	}

	txn, err := s.transactionRepo.Get(ctx, UID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewResourcesNotFound(err)
	}
	if checkAllowsOperation(txn, ownerID) {
		return api.TransactionResponse{}, mwerrors.NewResourceNotAllowedWithOnlyUsername(processor.GetUsername(ctx))
	}

	if input.Name != "" {
		txn.Name = input.Name
	}
	if input.Message != "" {
		txn.Status = input.Message
	}
	if input.Status != "" {
		txn.Status = input.Status
	}
	if input.ShowType != "" {
		txn.ShowType = input.ShowType
	}

	if err := s.transactionRepo.Update(ctx, txn); err != nil {
		return api.TransactionResponse{}, mwerrors.NewServerError(err)
	}
	resp, err := s.ConvertToApiTransaction(ctx, txn, false)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewServerError(err)
	}
	return resp, nil
}

// Delete deletes the transaction with the specified ID.
func (s service) Delete(ctx *gin.Context, ID string) (api.TransactionResponse, error) {
	UID, err := utils.Uint(ID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewIllegalArgumentError(err)
	}
	ownerID, err := utils.Uint(processor.GetUserID(ctx))
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewIllegalArgumentError(err)
	}

	txn, err := s.transactionRepo.Get(ctx, UID)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewServerError(err)
	}

	if checkAllowsOperation(txn, ownerID) {
		return api.TransactionResponse{}, mwerrors.NewResourceNotAllowedWithOnlyResourceID(processor.GetUsername(ctx), ownerID)
	}

	err = s.transactionRepo.Delete(ctx, txn)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewServerError(err)
	}
	resp, err := s.ConvertToApiTransaction(ctx, txn, false)
	if err != nil {
		return api.TransactionResponse{}, mwerrors.NewServerError(err)
	}
	return resp, nil
}

// Count returns the number of requester's transactions.
func (s service) Count(c *gin.Context) (string, int, error) {
	userID := processor.GetUserID(c)
	if userID == "" {
		return "Invalid", 0, mwerrors.NewMissingAuthToken()
	}
	return s.CountByUser(c, userID)
}

// CountByUser returns the number of user's transactions by user ID.
func (s service) CountByUser(ctx *gin.Context, ID string) (string, int, error) {
	userID := processor.GetUserID(ctx)
	if userID == "" {
		return "Invalid", 0, mwerrors.NewMissingAuthToken()
	}

	if userID == ID {
		return s.CountByUserWithShowType(ctx, ID, "Private")
	}

	requesterID, err := utils.Uint(userID)
	if err != nil {
		return "Invalid", 0, mwerrors.NewInvalidAuthToken(err)
	}
	requestUser, err := s.userRepo.Get(ctx, requesterID)
	if err != nil {
		return "Invalid", 0, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	otherID, err := utils.Uint(ID)
	if err != nil {
		return "Invalid", 0, mwerrors.NewIllegalArgumentError(err)
	}
	otherUser, err := s.userRepo.Get(ctx, otherID)
	if err != nil {
		return "Invalid", 0, mwerrors.NewResourceNotFoundWithPublicError(err)
	}

	showType := "Public"
	exists, err := s.friendRepo.HasRelationByBothUsers(ctx, requestUser, otherUser)
	if err != nil {
		return "Invalid", 0, mwerrors.NewResourceNotFoundWithPublicError(err)
	}
	if exists {
		showType = "Friend"
	}
	return s.CountByUserWithShowType(ctx, ID, showType)
}

// CountByUserWithShowType returns the number of user's transactions by user ID and show type.
func (s service) CountByUserWithShowType(c *gin.Context, ID string, showType string) (string, int, error) {
	ownerID, err := utils.Uint(ID)
	if err != nil {
		return "Invalid", 0, mwerrors.NewIllegalArgumentError(err)
	}
	cnt, err := s.transactionRepo.Count(c, ownerID, showType)
	return showType, cnt, err
}

// Query returns the transactions with the specified offset and limit.
func (s service) Query(c *gin.Context, ID, showType string, offset, limit int) ([]api.TransactionResponse, error) {
	userID := processor.GetUserID(c)
	if userID == "" {
		return []api.TransactionResponse{}, mwerrors.NewMissingAuthToken()
	}
	ownerID, err := utils.Uint(ID)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewIllegalInputErrorWithMessage(err.Error())
	}
	txns, err := s.transactionRepo.Query(c, offset, limit, ownerID, showType)
	if err != nil {
		return []api.TransactionResponse{}, mwerrors.NewResourcesNotFound(err)
	}
	resp, err := s.ConvertToApiTransactions(c, txns, showType != "Private")
	if err != nil {
		return []api.TransactionResponse{}, err
	}
	return resp, nil
}

func extractDeviceNameAndToken(devices []entity.Device) (string, []string) {
	var aggregatedIDs string
	var tokenList []string
	for idx, device := range devices {
		if idx == 0 {
			aggregatedIDs += strconv.Itoa(int(device.ID))
		} else {
			aggregatedIDs += "," + strconv.Itoa(int(device.ID))
		}
		tokenList = append(tokenList, device.DeviceToken)
	}
	return aggregatedIDs, tokenList
}

func (s service) ConvertToApiTransaction(c *gin.Context, txn entity.Transaction, isPrune bool) (api.TransactionResponse, error) {
	txns := []entity.Transaction{txn}
	res, err := s.ConvertToApiTransactions(c, txns, isPrune)
	if err != nil {
		return api.TransactionResponse{}, err
	}
	return res[0], nil
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

func checkAllowsOperation(transaction entity.Transaction, ownerID uint) bool {
	return transaction.SenderId != int(ownerID) && transaction.ReceiverId != int(ownerID)
}

func CreateTransactionMessage(requester entity.User, receiver entity.User, txn entity.Transaction) string {
	return fmt.Sprintf("Hi %s, %s sent you %f %s!", receiver.Username, requester.Username, txn.Amount, txn.Symbol)
}

func convert(txn entity.Transaction, sender, receiver entity.User, prune bool) api.TransactionResponse {
	if reflect.DeepEqual(sender, entity.User{}) {
		sender.Avatar = ""
		sender.Username = ""
	}
	if reflect.DeepEqual(receiver, entity.User{}) {
		receiver.Avatar = ""
		receiver.Username = ""
	}
	convertedTxn := api.Transaction{
		ID:               int(txn.ID),
		Name:             txn.Name,
		Status:           txn.Status,
		Amount:           "",
		Currency:         txn.Currency,
		Blockchain:       txn.Blockchain,
		Symbol:           txn.Symbol,
		SenderID:         txn.SenderId,
		SenderUrl:        sender.Avatar,
		SenderUsername:   sender.Username,
		SenderPubkey:     "",
		ReceiverID:       txn.ReceiverId,
		ReceiverUrl:      receiver.Avatar,
		ReceiverUsername: receiver.Username,
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
