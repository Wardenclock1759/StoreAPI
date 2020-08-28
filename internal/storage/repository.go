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

type RoleRepository interface {
	GrantRole(*model.UserRole) error
	RevokeRole(*model.UserRole) error
	GetRolesByID(uuid.UUID) (*model.UserRole, error)
}

type GameRepository interface {
	Create(*model.Game) error
	FindByName(string) (*model.Game, error)
	FindByID(uuid.UUID) (*model.Game, error)
}

type KeyRepository interface {
	Create(*model.Key) error
	FindByGame(uuid.UUID) (*[]string, error)
	FindByKey(string) (*model.Key, error)
}
