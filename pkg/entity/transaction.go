package entity

import "time"

type Transaction struct {
	Id   uint // primary_key
	Name string
	// Type            enum
	Status   string
	Amount   uint
	Currency string
	// Description     null
	SenderId       uint64
	SenderPubkey   uint64
	ReceiverId     uint64
	ReceiverPubkey uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Message        string
}
