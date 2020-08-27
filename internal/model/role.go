package model

import "github.com/google/uuid"

type Role string

const (
	RoleSeller Role = "seller"
)

type UserRole struct {
	Role Role      `json:"role"`
	ID   uuid.UUID `json:"user_id"`
}
