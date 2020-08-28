package sqlstorage

import (
	"database/sql"
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
	"github.com/google/uuid"
	"time"
)

type KeyRepository struct {
	storage *Storage
}

func (r *KeyRepository) Create(k *model.Key) error {
	if err := k.BeforeCreate(); err != nil {
		return err
	}

	r.storage.db.QueryRow(
		"INSERT INTO \"game_code\" (game_id, code, addedat) VALUES ($1, $2, $3)",
		k.ID,
		k.Key,
		k.AddedAt,
	)

	return nil
}

func (r *KeyRepository) Delete(k *model.Key) error {

	r.storage.db.QueryRow(
		"UPDATE \"game_code\" SET soldat = $1 WHERE game_id = $2 AND code = $3",
		time.Now(),
		k.ID,
		k.Key,
	)

	return nil
}

func (r *KeyRepository) FindByGame(id uuid.UUID) (*[]string, error) {
	var keys []string
	str := id.String()

	res, err := r.storage.db.Query(
		"SELECT (code) FROM \"game_code\" WHERE game_id = $1",
		str,
	)

	defer res.Close()
	var key string
	for res.Next() {
		err := res.Scan(&key)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err != nil {
		return nil, err
	}

	return &keys, nil
}

func (r *KeyRepository) FindByKey(key string) (*model.Key, error) {
	k := &model.Key{}
	if err := r.storage.db.QueryRow(
		"SELECT (game_id, code, addedat, removedat) FROM \"game_code\" WHERE code = $1",
		key,
	).Scan(&k.ID, &k.Key, &k.AddedAt, &k.SoldAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}

	return k, nil
}
