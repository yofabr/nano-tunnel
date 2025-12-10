package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show recent nano-tunnel CLI log hints",
	Long: `Print simple placeholder log lines to demonstrate where streaming
or historical logs could appear once implemented. The goal is to provide a
template for future log handling without performing any I/O.`,
	Example: "nano-tunnel logs",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now().Format(time.RFC3339)
		fmt.Printf("[%s] nano-tunnel: logs placeholder (hook up your logger here)\n", now)
		fmt.Printf("[%s] nano-tunnel: streaming would appear here\n", now)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
