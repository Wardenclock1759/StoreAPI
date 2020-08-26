package sqlstorage_test

import (
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
	"github.com/Wardenclock1759/StoreAPI/internal/storage/sqlstorage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstorage.TestDB(t, databaseURL)
	defer teardown("user")

	s := sqlstorage.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstorage.TestDB(t, databaseURL)
	defer teardown("user")

	s := sqlstorage.New(db)
	email := "test@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, storage.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByID(t *testing.T) {
	db, teardown := sqlstorage.TestDB(t, databaseURL)
	defer teardown("user")

	s := sqlstorage.New(db)

	u1 := model.TestUser(t)
	u1.ID = uuid.New()
	s.User().Create(u1)

	u2, err := s.User().FindByID(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
