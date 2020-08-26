package storage_test

import (
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("user")

	u, err := s.User().Create(model.TestUser(t))
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("user")

	email := "test@example.org"
	_, err := s.User().FindByEmail(email)
	assert.Error(t, err)

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
