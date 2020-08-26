package apiserver

import (
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	router  *mux.Router
	logger  *logrus.Logger
	storage storage.Storage
}

func newServer(storage storage.Storage) *server {
	s := &server{
		router:  mux.NewRouter(),
		logger:  logrus.New(),
		storage: storage,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {

}
