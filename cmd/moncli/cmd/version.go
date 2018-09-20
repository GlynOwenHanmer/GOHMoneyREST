package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// to be changed using ldflags with the go build command
var version = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("generate a bash completion script for %s", appName),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
