package cmd

import (
	"Coins/APP"
	"errors"
	"os"
)
import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "go-coin",
	Short:        "go-coin",
	SilenceUsage: true,
	Long:         `go-coin`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("nothing command")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
func init() {
	rootCmd.AddCommand(APP.StartCmd)
	rootCmd.AddCommand(APP.Account)
}
