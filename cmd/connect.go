package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [remote-host]",
	Short: "Simulate connecting to a remote nano-tunnel endpoint",
	Long: `Demonstrate the shape of a connect workflow. This command does not
open a real connection; it simply echoes the target host and prints a mock
connection lifecycle. Replace the placeholders with the real dial logic when
ready.`,
	Example: "nano-tunnel connect nano-tunnel.onrender.com",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			_ = cmd.Help()
			return
		}

		target := args[0]
		fmt.Printf("Connecting to %s ...\n", target)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Connected to %s (placeholder)\n", target)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
