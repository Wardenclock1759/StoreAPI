package model

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID          uuid.UUID `json:"id"`
	Time        time.Time `json:"time"`
	Card        string    `json:"card,omitempty"`
	GameName    string    `json:"game_name"`
	UserEmail   string    `json:"user_email"`
	SellerEmail string    `json:"seller_email"`
	Total       int       `json:"total"`
	Code        string    `json:"code"`
}

func (p *Payment) PostCreate() error {
	p.ID = uuid.New()
	return nil
}
