package entity

import (
	"fmt"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Transaction struct {
	Id             uint   `json:"id"` // string of hex? then use "hexadecimal"
	Name           string `json:"name" validate:"required"`
	Status         string
	Amount         uint      `json:"amount" validate:"required,uint"`
	Currency       string    `json:"currency" validate:"required,iso4217"` //currency code
	SenderId       uint64    `json:"sender_id" validate:"uuid"`
	SenderPubkey   uint64    `json:"sender_pk" validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	ReceiverId     uint64    `json:"receiver_id" validate:"uuid"`
	ReceiverPubkey uint64    `json:"receiver_pk" validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	CreatedAt      time.Time `json:"creat_at" validate:"required,datetime"`
	UpdatedAt      time.Time `json:"update_at" validate:"required,datetime"`
	// message should be less than 200 characters
	Message string `json:"message" validate:"ls=200"`
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
