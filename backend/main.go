package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"metavida/backend/config"
	"metavida/backend/core"
	"metavida/backend/db"
	"metavida/backend/security"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "deploy_tables":
			deployTablesCommand()
			return
		case "insert_admin":
			insertAdminCommand()
			return
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	core.SetAuthSecret(cfg.SecretPhrase)
	db.Configure(cfg)
	makeAppHandlers()

	if cfg.CloudProvider == "cloudflare" && !cfg.HasCloudflareD1Credentials() {
		log.Println("Cloudflare D1 credentials are missing; skipping database initialization")
	} else {
		if err := db.DeployTables(cfg, MakeDeploySchemas()...); err != nil {
			log.Fatal(err)
		}
		if err := security.EnsureBootstrapAdmin(cfg); err != nil {
			log.Fatal(err)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", cors(cfg, http.HandlerFunc(LocalHandler)))

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Printf("Metavida backend listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func deployTablesCommand() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if !cfg.HasCloudflareD1Credentials() {
		log.Fatal("missing Cloudflare D1 credentials in credentials.json or METAVIDA_CREDENTIALS")
	}
	schemas := MakeDeploySchemas()
	if err := db.DeployTables(cfg, schemas...); err != nil {
		log.Fatal(err)
	}
	log.Printf("deployed %d table schemas to Cloudflare D1", len(schemas))
}

func insertAdminCommand() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if !cfg.HasCloudflareD1Credentials() {
		log.Fatal("missing Cloudflare D1 credentials in credentials.json or METAVIDA_CREDENTIALS")
	}
	db.Configure(cfg)
	if err := security.EnsureBootstrapAdmin(cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("bootstrap admin user ensured for %s", cfg.AdminEmail)
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
