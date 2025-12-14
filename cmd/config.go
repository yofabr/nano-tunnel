package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [path]",
	Short: "Preview and validate a nano-tunnel config file",
	Long: `Read a nano-tunnel JSON config, normalize the remote host, and print
the result. This is a dry-run helper for verifying configuration before
starting the tunnel.`,
	Example: "nano-tunnel config ./your_config_file.json",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter config file name:")
		scanner.Scan()
		filename := scanner.Text()
		if filename == "" {
			filename = "./config.json"
		}
		if _, err := os.Stat(filename); err == nil {
			fmt.Println("Config file already exists")
			return
		}

		fmt.Println("Enter Remote URL (default: nano-tunnel.onrender.com)")
		scanner.Scan()
		url := scanner.Text()
		if url == "" {
			url = "nano-tunnel.onrender.com"
		}

		payload := map[string]any{
			"remote_url": url,
		}

		fmt.Println("Enter Config Path (default: ./)")
		scanner.Scan()
		path := scanner.Text()
		if path == "" {
			path = "./"
		}

		out, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Println("error rendering config preview:", err)
			return
		}

		filename = strings.TrimSuffix(filename, ".json")
		filename = filename + ".json"
		os.WriteFile(filename, out, 0644)
		fmt.Println("Config file created:", filename)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
