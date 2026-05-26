package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"metavida/backend/config"
	"metavida/backend/core"
	"metavida/backend/db"
	"metavida/backend/security"
)

type Server struct {
	cfg     config.Config
	handler http.Handler
}

func New(cfg config.Config) (*Server, error) {
	db.Configure(cfg)
	s := &Server{cfg: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/clients", s.clients)
	s.handler = cors(cfg, mux)
	return s, nil
}

func (s *Server) Init() error {
	if s.cfg.CloudProvider == "cloudflare" && !s.cfg.HasCloudflareD1Credentials() {
		log.Println("Cloudflare D1 credentials are missing; skipping database initialization")
		return nil
	}
	return db.DeployTables(s.cfg, db.Table[security.User]())
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.HTTPAddr, s.handler)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	core.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) clients(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createClient(w, r)
	case http.MethodGet:
		s.listClients(w)
	default:
		core.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

type clientInput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func (s *Server) createClient(w http.ResponseWriter, r *http.Request) {
	var input clientInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		core.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Message = strings.TrimSpace(input.Message)
	if input.Name == "" {
		core.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	now := core.SUnixTime()
	user := security.User{
		ID:      core.SUnixTime(),
		Name:    input.Name,
		Email:   input.Email,
		Phone:   input.Phone,
		Message: input.Message,
		Type:    1,
		Status:  1,
		Created: now,
		Updated: now,
	}
	if err := db.InsertOne(user); err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	core.JSON(w, http.StatusCreated, user)
}

func (s *Server) listClients(w http.ResponseWriter) {
	rows := []security.User{}
	query := db.Query(&rows)
	query.Type.Equals(int8(1)).Limit(50)
	if err := query.Exec(); err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	core.JSON(w, http.StatusOK, rows)
}

func cors(cfg config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := cfg.FrontendOrigin
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
