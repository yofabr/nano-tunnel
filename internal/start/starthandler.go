package start

import (
	"encoding/json"
	"os"
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

func NewListener(path string) (*AppData, error) {
	con, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	appData := AppData{}
	json.Unmarshal(con, &appData)
	return &appData, nil
}
