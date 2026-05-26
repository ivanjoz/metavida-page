package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"metavida/backend/config"
)

type d1Request struct {
	SQL    string `json:"sql"`
	Params []any  `json:"params,omitempty"`
}

type d1Response struct {
	Success bool `json:"success"`
	Result  []struct {
		Results []map[string]any `json:"results"`
		Success bool             `json:"success"`
		Error   string           `json:"error"`
	} `json:"result"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var (
	d1Mu     sync.RWMutex
	d1Cfg    config.Config
	d1Client = &http.Client{Timeout: 20 * time.Second}
)

func Configure(cfg config.Config) {
	d1Mu.Lock()
	defer d1Mu.Unlock()
	d1Cfg = cfg
	d1Client = &http.Client{Timeout: 20 * time.Second}
}

func execD1(sql string, params []any) (d1Response, error) {
	d1Mu.RLock()
	cfg := d1Cfg
	client := d1Client
	d1Mu.RUnlock()

	var empty d1Response
	if cfg.CloudflareAccountID == "" || cfg.CloudflareAPIToken == "" || cfg.CloudflareDatabaseID == "" {
		return empty, errors.New("missing Cloudflare D1 credentials; call db.Configure or db.DeployTables with credentials")
	}

	payload, err := json.Marshal(d1Request{SQL: sql, Params: params})
	if err != nil {
		return empty, err
	}
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/query",
		cfg.CloudflareAccountID, cfg.CloudflareDatabaseID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return empty, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.CloudflareAPIToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return empty, err
	}
	defer res.Body.Close()

	var out d1Response
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return empty, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 || !out.Success {
		if len(out.Errors) > 0 {
			return empty, errors.New(out.Errors[0].Message)
		}
		return empty, fmt.Errorf("cloudflare d1 query failed with status %d", res.StatusCode)
	}
	for _, result := range out.Result {
		if !result.Success {
			return empty, errors.New(result.Error)
		}
	}
	return out, nil
}
