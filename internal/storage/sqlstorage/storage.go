package sqlstorage

import (
	"database/sql"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db             *sql.DB
	userRepository *UserRepository
	roleRepository *RoleRepository
	gameRepository *GameRepository
	keyRepository  *KeyRepository
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

func (s *Storage) Role() storage.RoleRepository {
	if s.roleRepository != nil {
		return s.roleRepository
	}

	s.roleRepository = &RoleRepository{
		storage: s,
	}

	return s.roleRepository
}

func (s *Storage) Game() storage.GameRepository {
	if s.gameRepository != nil {
		return s.gameRepository
	}

	s.gameRepository = &GameRepository{
		storage: s,
	}

	return s.gameRepository
}

func (s *Storage) Key() storage.KeyRepository {
	if s.keyRepository != nil {
		return s.keyRepository
	}

	s.keyRepository = &KeyRepository{
		storage: s,
	}

	return s.keyRepository
}
