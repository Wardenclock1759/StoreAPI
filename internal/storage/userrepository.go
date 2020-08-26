package storage

import "github.com/Wardenclock1759/StoreAPI/internal/model"

type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}

	r.storage.db.QueryRow(
		"INSERT INTO \"user\" (user_id, email, encrypted_password) VALUES ($1, $2, $3)",
		u.ID,
		u.Email,
		u.EncryptedPassword,
	)

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.storage.db.QueryRow(
		"SELECT user_id, email, encrypted_password FROM \"user\" WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
		return nil, err
	}

	return u, nil
}
