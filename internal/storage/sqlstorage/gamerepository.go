package sqlstorage

import (
	"database/sql"
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
	"github.com/google/uuid"
)

type GameRepository struct {
	storage *Storage
}

func (r *GameRepository) Create(g *model.Game) error {
	if err := g.Validate(); err != nil {
		return err
	}

	if err := g.BeforeCreate(); err != nil {
		return err
	}

	r.storage.db.QueryRow(
		"INSERT INTO \"game\" (game_id, user_id, name, price) VALUES ($1, $2, $3, $4)",
		g.ID,
		g.User,
		g.Name,
		g.Price,
	)

	return nil
}

func (r *GameRepository) FindByID(id uuid.UUID) (*model.Game, error) {
	g := &model.Game{}
	if err := r.storage.db.QueryRow(
		"SELECT game_id, user_id, name, price FROM \"game\" WHERE game_id = $1",
		id,
	).Scan(&g.ID, &g.User, &g.Name, &g.Price); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}

	return g, nil
}

func (r *GameRepository) FindByName(name string) (*model.Game, error) {
	g := &model.Game{}
	if err := r.storage.db.QueryRow(
		"SELECT game_id, user_id, name, price FROM \"game\" WHERE name = $1",
		name,
	).Scan(&g.ID, &g.User, &g.Name, &g.Price); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}

	return g, nil
}
