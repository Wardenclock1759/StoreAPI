package model

import (
	"github.com/google/uuid"
	"time"
)

type PaymentSession struct {
	ID          uuid.UUID `json:"id"`
	Time        time.Time `json:"time"`
	Card        string    `json:"card,omitempty"`
	GameName    string    `json:"game_name"`
	UserEmail   string    `json:"user_email"`
	SellerEmail string    `json:"seller_email"`
}
