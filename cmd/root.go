/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nano-tunnel",
	Short: "Forward your local ports to the internet via a lightweight tunnel",
	Long: `Nano-tunnel is a tiny Go CLI that keeps a persistent WebSocket
connection to a remote Nano-tunnel server so you can securely forward HTTP
requests to your local machine from anywhere.

Typical flow:
  1) Create a config file with your remote server:
       { "remote_url": "nano-tunnel.onrender.com" }
  2) Start the tunnel with your config:
       nano-tunnel start ./your_config_file.json
  3) Copy the printed Client ID into the hosted Nano-tunnel UI and forward
     requests to the local port you specify.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
