package storage_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=store_api_db_test user=User password=123456 sslmode=disable"
	}

	os.Exit(m.Run())
}
