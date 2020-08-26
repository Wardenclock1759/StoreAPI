package model

import (
	"github.com/google/uuid"
	"testing"
)

func TestUser(t *testing.T) *User {
	return &User{
		ID:       uuid.New(),
		Email:    "test@example.org",
		Password: "123456789",
	}
}
