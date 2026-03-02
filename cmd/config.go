package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	configPath      string
	configURL       string
	configOverwrite bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate a nano-tunnel config file",
	Long: `Create a nano-tunnel JSON config file with the specified remote URL.

Example:
  nano-tunnel config -u nano-tunnel.onrender.com -o ./config.json`,
	Example: "nano-tunnel config -u nano-tunnel.onrender.com -o ./config.json",
	Run: func(cmd *cobra.Command, args []string) {
		filename := configPath
		if filename == "" {
			filename = "./config.json"
		}

		filename = strings.TrimSuffix(filename, ".json") + ".json"

		if !configOverwrite {
			if _, err := os.Stat(filename); err == nil {
				fmt.Printf("Config file %s already exists. Use -f to overwrite.\n", filename)
				return
			}
		}

		if configURL == "" {
			configURL = "nano-tunnel.onrender.com"
		}

		payload := map[string]any{
			"remote_url": configURL,
		}

		out, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Println("error rendering config:", err)
			return
		}

		if err := os.WriteFile(filename, out, 0644); err != nil {
			fmt.Println("error writing config file:", err)
			return
		}

		fmt.Println("Config file created:", filename)
	},
}

func init() {
	configCmd.Flags().StringVarP(&configPath, "output", "o", "", "Output config file path (default: ./config.json)")
	configCmd.Flags().StringVarP(&configURL, "url", "u", "", "Remote server URL (default: nano-tunnel.onrender.com)")
	configCmd.Flags().BoolVarP(&configOverwrite, "force", "f", false, "Overwrite existing config file")
	rootCmd.AddCommand(configCmd)
}
