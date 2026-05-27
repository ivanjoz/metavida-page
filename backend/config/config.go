package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	AppEnv               string `json:"APP_ENV"`
	HTTPAddr             string `json:"HTTP_ADDR"`
	FrontendOrigin       string `json:"FRONTEND_ORIGIN"`
	CloudProvider        string `json:"CLOUD_PROVIDER"`
	CloudflareAccountID  string `json:"CLOUDFLARE_ACCOUNT"`
	CloudflareAPIToken   string `json:"CLOUDFLARE_TOKEN"`
	CloudflareDatabaseID string `json:"CLOUDFLARE_DATABASE_ID"`
	AdminEmail           string `json:"ADMIN_EMAIL"`
	AdminPassword        string `json:"ADMIN_PASSWORD"`
	SecretPhrase         string `json:"SECRET_PHRASE"`
	AWSUserKey           string `json:"AWS_USER_KEY"`
	AWSSecretKey         string `json:"AWS_SECRET_KEY"`
	AWSRegion            string `json:"AWS_REGION"`
	AWSSESSender         string `json:"AWS_SES_SENDER"`
	AWSSESRecipient      string `json:"AWS_SES_RECIPIENT"`
	SenderEmail          string `json:"SENDER_EMAIL"`
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:         "local",
		HTTPAddr:       ":3591",
		FrontendOrigin: "http://localhost:3571",
		CloudProvider:  "cloudflare",
	}

	path := os.Getenv("METAVIDA_CREDENTIALS")
	if path == "" {
		path = filepath.Join("credentials.json")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read credentials: %w", err)
	}

	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return cfg, fmt.Errorf("parse credentials: %w", err)
	}

	return cfg, nil
}

func (cfg Config) HasCloudflareD1Credentials() bool {
	return cfg.CloudflareAccountID != "" &&
		cfg.CloudflareAPIToken != "" &&
		cfg.CloudflareDatabaseID != ""
}
