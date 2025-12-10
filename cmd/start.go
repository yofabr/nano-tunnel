package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/yofabr/nano-tunnel/internal/forward"
	"github.com/yofabr/nano-tunnel/internal/logger"
	"github.com/yofabr/nano-tunnel/internal/start"
)

type WsData struct {
	LocalPort string                 `json:"local_port,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Method    string                 `json:"method,omitempty"`
	Headers   map[string]string      `json:"headers,omitempty"`
	Body      map[string]interface{} `json:"body,omitempty"`
}

type Message struct {
	Event    string `json:"event"`
	ClientID string `json:"clientID,omitempty"`
	Message  string `json:"message,omitempty"`
	Data     WsData `json:"data,omitempty"`
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts Polling local port to the internet",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("start called")
		if len(args) > 1 || len(args) == 0 {
			fmt.Println("ERROR: Enter valid path to the config file")
			return
		}

		listener, err := start.NewListener(args[0])
		if err != nil {
			log.Fatal("Error while reading:", err)
		}
		u := url.URL{Scheme: "ws", Host: listener.RemoteUrl, Path: "/ws"}
		fmt.Println("Connecting to", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

		if err != nil {
			log.Fatal("Dial error:", err)
		}
		defer c.Close()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			var m Message
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
				fmt.Println(m)
				url := fmt.Sprintf("http://localhost:%s%s", m.Data.LocalPort, m.Data.Path)
				forward.FetchResource(c, url, m.Data.Method, m.Data.Headers, m.Data.Body)

			default:
				log.Println("Unknown event:", m.Event)
			}
		}
		fmt.Println("Listener", *listener)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
