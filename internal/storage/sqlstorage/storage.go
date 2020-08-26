package sqlstorage

import (
	"database/sql"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) User() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		storage: s,
	}

	return s.userRepository
}
