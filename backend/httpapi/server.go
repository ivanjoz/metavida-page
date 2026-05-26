package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"metavida/backend/config"
	"metavida/backend/core"
	"metavida/backend/db"
	"metavida/backend/domain/leads"
)

type Server struct {
	cfg     config.Config
	leadDB  db.ORM[leads.Lead]
	handler http.Handler
}

func New(cfg config.Config) (*Server, error) {
	leadDB, err := db.NewORM[leads.Lead](cfg, "leads")
	if err != nil {
		return nil, err
	}

	s := &Server{cfg: cfg, leadDB: leadDB}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/leads", s.leads)
	s.handler = cors(cfg, mux)
	return s, nil
}

func (s *Server) Init() error {
	if s.cfg.CloudProvider == "cloudflare" && !s.cfg.HasCloudflareD1Credentials() {
		log.Println("Cloudflare D1 credentials are missing; skipping database initialization")
		return nil
	}
	return s.leadDB.Init()
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.cfg.HTTPAddr, s.handler)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	core.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) leads(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createLead(w, r)
	case http.MethodGet:
		s.listLeads(w)
	default:
		core.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

type leadInput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (s *Server) createLead(w http.ResponseWriter, r *http.Request) {
	var input leadInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		core.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Message = strings.TrimSpace(input.Message)
	if input.Name == "" {
		core.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	lead := leads.New(input.Name, input.Email, input.Message, uuid.NewString())
	if err := s.leadDB.Insert(lead); err != nil {
		core.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	core.JSON(w, http.StatusCreated, lead)
}

func (s *Server) listLeads(w http.ResponseWriter) {
	rows, err := s.leadDB.Select().Limit(50).Exec()
	if err != nil {
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
