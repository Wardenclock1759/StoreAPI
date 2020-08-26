package sqlstorage

import (
	"database/sql"
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
)

type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	r.storage.db.QueryRow(
		"INSERT INTO \"user\" (user_id, email, encrypted_password) VALUES ($1, $2, $3)",
		u.ID,
		u.Email,
		u.EncryptedPassword,
	)

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow(
		"SELECT user_id, email, encrypted_password FROM \"user\" WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}
