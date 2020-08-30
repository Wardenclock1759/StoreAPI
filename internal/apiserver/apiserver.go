package apiserver

import (
	"database/sql"
	"github.com/Wardenclock1759/StoreAPI/internal/storage/sqlstorage"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

func Start(config *Config) error {
	url := os.Getenv("DATABASE_URL")
	db, err := newDB(url)
	if err != nil {
		return err
	}

	defer db.Close()
	storage := sqlstorage.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(storage, sessionStore)

	var dynamicPort string
	port := os.Getenv("PORT")
	if port != "" {
		dynamicPort = port
	} else {
		dynamicPort = config.BindAddress
	}
	return http.ListenAndServe(dynamicPort, s)
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
