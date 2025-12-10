package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display environment and binary information",
	Long: `Print a snapshot of your nano-tunnel CLI environment, including
runtime details, Go version, platform metadata, and the current binary path.

This is useful when filing issues or confirming that your build matches the
expected platform requirements. No network calls are made by this command.`,
	Example: "nano-tunnel info",
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]any{
			"timestamp":  time.Now().Format(time.RFC3339),
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
		}

		if exe, err := os.Executable(); err == nil {
			payload["binary_path"] = exe
		}

		out, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Println("error building info payload:", err)
			return
		}

		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
