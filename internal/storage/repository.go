package storage

import "github.com/Wardenclock1759/StoreAPI/internal/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}
