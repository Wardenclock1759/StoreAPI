package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/Wardenclock1759/StoreAPI/internal/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrFailedToGenerateToken    = errors.New("failed to generate token string")
	ErrFailedDecodeToken        = errors.New("failed to decode token")
	ErrForbidden                = errors.New("forbidden to proceed")
	signedKey                   = []byte(os.Getenv("JWT_KEY"))
)

type ctxKey int8

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
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/user/sign-up", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/user/sign-in", s.handleJWTCreate()).Methods("POST")

	role := s.router.PathPrefix("/user/role").Subrouter()
	role.Use(s.authorisedUser)
	role.HandleFunc("/grant-role", s.handleGrantRole()).Methods("POST")
	role.HandleFunc("/revoke-role", s.handleRevokeRole()).Methods("POST")

	store := s.router.PathPrefix("/store").Subrouter()
	store.Use(s.authorisedUser)

	seller := store.PathPrefix("/publisher").Subrouter()
	seller.Use(s.authorisedSeller)
	seller.HandleFunc("/game", s.handleGameCreate()).Methods("POST")
	seller.HandleFunc("/key", s.handleKeyCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authorisedUser)
	private.HandleFunc("/roles", s.handleGetRole()).Methods("GET")
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleGetRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		role, err := s.storage.Role().GetRolesByID(u.ID)
		if err != nil {

		}
		s.respond(w, r, http.StatusOK, role)
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"Remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start))
	})
}

func (s *server) authorisedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := verifyToken(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrFailedDecodeToken)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			s.error(w, r, http.StatusInternalServerError, ErrFailedDecodeToken)
			return
		}
		id, valid := claims["user_id"].(string)

		if !valid {
			s.error(w, r, http.StatusInternalServerError, ErrFailedDecodeToken)
			return
		}
		id2, err2 := uuid.Parse(id)
		u, err := s.storage.User().FindByID(id2)
		if err != nil || err2 != nil {
			s.error(w, r, http.StatusForbidden, ErrForbidden)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) authorisedSeller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := verifyToken(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrFailedDecodeToken)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			s.error(w, r, http.StatusInternalServerError, ErrFailedDecodeToken)
			return
		}
		id, valid := claims["user_id"].(string)

		if !valid {
			s.error(w, r, http.StatusInternalServerError, ErrFailedDecodeToken)
			return
		}
		id2, err2 := uuid.Parse(id)
		u, err := s.storage.User().FindByID(id2)
		if err != nil || err2 != nil {
			s.error(w, r, http.StatusInternalServerError, ErrFailedDecodeToken)
			return
		}

		role, err := s.storage.Role().GetRolesByID(id2)
		if role == nil || err != nil {
			s.error(w, r, http.StatusForbidden, ErrForbidden)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleUserCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.storage.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	if r.Header["Token"] != nil {
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrFailedDecodeToken
			}
			return signedKey, nil
		})

		if err != nil {
			return nil, err
		}

		return token, nil
	}
	return nil, ErrFailedDecodeToken
}

func (s *server) handleJWTCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.storage.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		jwtString, err := GenerateJWT(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrFailedToGenerateToken)
			return
		}

		w.Header().Set("Token", jwtString)

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGrantRole() http.HandlerFunc {
	type request struct {
		ID   uuid.UUID  `json:"user_id"`
		Role model.Role `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		role := &model.UserRole{
			ID:   req.ID,
			Role: req.Role,
		}
		if err := s.storage.Role().GrantRole(role); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, role)
	}
}

func (s *server) handleRevokeRole() http.HandlerFunc {
	type request struct {
		ID   uuid.UUID  `json:"user_id"`
		Role model.Role `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		role := &model.UserRole{
			ID:   req.ID,
			Role: req.Role,
		}
		if err := s.storage.Role().RevokeRole(role); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, role)
	}
}

func (s *server) handleGameCreate() http.HandlerFunc {
	type request struct {
		Name  string    `json:"name"`
		Price string    `json:"price"`
		User  uuid.UUID `json:"user_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		g := &model.Game{
			Name:  req.Name,
			Price: req.Price,
			User:  req.User.String(),
		}
		if err := s.storage.Game().Create(g); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, g)
	}
}

func (s *server) handleKeyCreate() http.HandlerFunc {
	type request struct {
		ID  uuid.UUID `json:"game_id"`
		Key string    `json:"code"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		k := &model.Key{
			ID:  req.ID,
			Key: req.Key,
		}
		if err := s.storage.Key().Create(k); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, k)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func GenerateJWT(u *model.User) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := at.SignedString(signedKey)

	if err != nil {
		return "", ErrFailedToGenerateToken
	}

	return tokenString, nil
}