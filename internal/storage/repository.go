package storage

import (
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	FindByID(uuid.UUID) (*model.User, error)
}
