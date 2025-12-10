package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yofabr/nano-tunnel/internal/start"
)

var statusCmd = &cobra.Command{
	Use:   "status [config.json]",
	Short: "Show the current tunnel status from a config file",
	Long: `Load the provided config file and print a lightweight status summary.
This does not open a WebSocket; it only validates and echoes what would be used
by the tunnel. Use it to confirm your remote host and local port before running
the start command.`,
	Example: "nano-tunnel status ./your_config_file.json",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}

		cfg, err := start.NewListener(args[0])
		if err != nil {
			fmt.Println("failed to read config:", err)
			return
		}

		summary := map[string]any{
			"remote_host": cfg.RemoteUrl,
			"local_port":  cfg.LocalPort,
			"checked_at":  time.Now().Format(time.RFC3339),
		}

		out, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			fmt.Println("failed to render status:", err)
			return
		}

		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
