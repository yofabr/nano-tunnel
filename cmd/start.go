package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/yofabr/nano-tunnel/internal/forward"
	"github.com/yofabr/nano-tunnel/internal/logger"
	"github.com/yofabr/nano-tunnel/internal/start"
	"github.com/yofabr/nano-tunnel/internal/types"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Connects to the Nano-tunnel server using a JSON config file",
	Long: `Start a WebSocket connection to the Nano-tunnel server so the hosted
UI can forward HTTP requests to your local machine.

Example config file (your_config_file.json):
  {
    "remote_url": "nano-tunnel.onrender.com"
  }

Run the tunnel:
  nano-tunnel start ./your_config_file.json

Keep this command running, copy the printed Client ID into the hosted UI,
then enter the local port and API path you want to forward.`,
	Example: "nano-tunnel start ./your_config_file.json",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 || len(args) == 0 {
			_ = cmd.Help()
			return
		}

		configPath := args[0]
		if _, err := os.Stat(configPath); err != nil {
			log.Fatalf("unable to read config file %s: %v", configPath, err)
		}

		listener, err := start.NewListener(configPath)
		if err != nil {
			log.Fatal("Error while reading config:", err)
		}

		u := url.URL{Scheme: "wss", Host: listener.RemoteUrl, Path: "/ws"}
		fmt.Printf("Connecting to %s (config: %s)\n", u.String(), filepath.Clean(configPath))

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

		if err != nil {
			log.Fatal("Dial error:", err)
		}

		defer logger.DisconnectLogger()
		defer c.Close()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			var m types.Message
			if err := json.Unmarshal(msg, &m); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}

			switch m.Event {
			case "welcome":
				logger.WelcomLogger(m.ClientID)
			case "broad":
				log.Println("Broadcast message:", m.Message)

			case "forward":
				forwardURL := fmt.Sprintf("http://localhost:%s%s", m.Data.LocalPort, m.Data.Path)
				log.Printf("Forwarding request for client %s to %s", m.ClientID, forwardURL)
				forward.FetchResource(c, m.ClientID, forwardURL, m.Data.Method, m.Data.Headers, m.Data.Body)

			default:
				log.Println("Unknown event:", m.Event)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
