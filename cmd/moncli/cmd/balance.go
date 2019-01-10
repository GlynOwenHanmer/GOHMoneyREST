package cmd

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "interact with a balance",
}

var balanceDeleteCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "delete a balance",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		id, err := parseID(args[0])
		if err != nil {
			return errors.Wrap(err, "parsing balance ID")
		}
		c, err := newClient()
		if err != nil {
			return errors.Wrap(err, "creating new client")
		}
		return c.DeleteBalance(uint(id))
	},
}

func init() {
	err := viper.BindPFlags(balanceCmd.PersistentFlags())
	if err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
	rootCmd.AddCommand(balanceCmd)

	for _, c := range []*cobra.Command{
		balanceDeleteCmd,
	} {
		err := viper.BindPFlags(c.Flags())
		if err != nil {
			log.Fatal(errors.Wrap(err, "binding pflags"))
		}
		balanceCmd.AddCommand(c)
	}
}
