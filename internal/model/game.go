package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type Game struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	User  uuid.UUID `json:"user"`
	Price string    `json:"price"`
}

func (g *Game) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(&g.Name, validation.Required),
		validation.Field(&g.Price, validation.Required),
		validation.Field(&g.ID, validation.Required, is.UUID),
		validation.Field(&g.User, validation.Required, is.UUID),
	)
}

func (g *Game) BeforeCreate() error {
	g.ID = uuid.New()

	return nil
}
