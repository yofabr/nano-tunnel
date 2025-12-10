package start

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type AppData struct {
	RemoteUrl string   `json:"remote_url"`
	Path      int      `json:"path"`
	LocalPort string   `json:"local_port"`
	Cors      []string `json:"cors"`
}

func (a *AppData) start_polling() {
	// Initialize a web socket here to forward incoming requests...
}

func (a *AppData) Validate() error {
	if a == nil {
		return errors.New("config not loaded")
	}

	normalized, err := normalizeRemoteURL(a.RemoteUrl)
	if err != nil {
		return err
	}

	a.RemoteUrl = normalized
	return nil
}

func NewListener(path string) (*AppData, error) {
	con, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	appData := AppData{}
	if err := json.Unmarshal(con, &appData); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if err := appData.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &appData, nil
}

func normalizeRemoteURL(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", errors.New("remote_url is required")
	}

	raw = strings.TrimSpace(raw)

	for _, prefix := range []string{"wss://", "ws://", "https://", "http://"} {
		raw = strings.TrimPrefix(raw, prefix)
	}

	raw = strings.TrimSuffix(raw, "/")

	if strings.Contains(raw, "/") {
		return "", fmt.Errorf("remote_url should only be a host (got %s)", raw)
	}

	return raw, nil
}
