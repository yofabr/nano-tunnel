package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yofabr/nano-tunnel/internal/start"
)

var configCmd = &cobra.Command{
	Use:   "config [path]",
	Short: "Preview and validate a nano-tunnel config file",
	Long: `Read a nano-tunnel JSON config, normalize the remote host, and print
the result. This is a dry-run helper for verifying configuration before
starting the tunnel.`,
	Example: "nano-tunnel config ./your_config_file.json",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}

		cfgPath := args[0]
		cfg, err := start.NewListener(cfgPath)
		if err != nil {
			fmt.Println("error reading config:", err)
			return
		}

		payload := map[string]any{
			"path":       filepath.Clean(cfgPath),
			"remote_url": cfg.RemoteUrl,
			"local_port": cfg.LocalPort,
		}

		out, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Println("error rendering config preview:", err)
			return
		}

		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
