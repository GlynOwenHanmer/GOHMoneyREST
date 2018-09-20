package versioncmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func New(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show the version of this application",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(version)
			return nil
		},
	}
}
