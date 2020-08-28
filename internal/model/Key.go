package model

import (
	"github.com/google/uuid"
	"time"
)

type Key struct {
	ID       uuid.UUID `json:"game_id"`
	Key      string    `json:"key"`
	AddedAt  time.Time `json:"added_at"`
	SoldAt   time.Time `json:"sold_at,omitempty"`
	BoughtBy uuid.UUID `json:"bought_by,omitempty"`
}

func (k *Key) BeforeCreate() error {
	k.AddedAt = time.Now()

	return nil
}
