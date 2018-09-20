package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// to be changed using ldflags with the go build command
var version = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version of this application",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version)
		return nil
	},
}

func init() {
	cmdDBServe.AddCommand(versionCmd)
}
