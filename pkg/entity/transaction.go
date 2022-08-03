package entity

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

// Transaction transaction
//
// swagger:model transaction
type Transaction struct {
	gorm.Model            // adds ID, created_at etc.
	Name           string `json:"name"        validate:"required"`
	Status         string `json:"status"`
	Amount         uint   `json:"amount"      validate:"required,uint"`
	Currency       string `json:"currency"    validate:"required,iso4217"` //currency code
	SenderId       uint   `json:"sender_id"   validate:"uuid"`
	SenderPubkey   uint64 `json:"sender_pk"   validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	ReceiverId     uint   `json:"receiver_id" validate:"uuid"`
	ReceiverPubkey uint64 `json:"receiver_pk" validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	// message should be less than 200 characters
	Message string `json:"message"     validate:"ls=200"`
}

// Validate validates the UpdateTransactionRequest fields.
func (m Transaction) Validate() error {
	validate := validator.New()
	err := validate.StructExcept(m, "Status")
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
