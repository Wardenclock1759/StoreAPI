package apiserver

import (
	"database/sql"
	"github.com/Wardenclock1759/StoreAPI/internal/storage/sqlstorage"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	storage := sqlstorage.New(db)
	s := newServer(storage)

	return http.ListenAndServe(config.BindAddress, s)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
