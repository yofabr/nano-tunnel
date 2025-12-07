/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/yofabr/nano_tunnel/internal/start"
)

type Message struct {
	Event    string `json:"event"`
	ClientID string `json:"clientID,omitempty"`
	Message  string `json:"message,omitempty"`
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
		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
		log.Println("Connecting to", u.String())

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
				log.Println("Connected! ClientID:", m.ClientID)
			case "broad":
				log.Println("Broadcast message:", m.Message)
			default:
				log.Println("Unknown event:", m.Event)
			}
		}
		fmt.Println("Listener", *listener)

		// for {
		// 	time.Sleep(1 * time.Second)
		// }
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
