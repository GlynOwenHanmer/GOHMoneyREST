package cmd

import (
	"log"
	"os"

	"github.com/glynternet/mon/pkg/storage"
	"github.com/glynternet/mon/pkg/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "interact with a balance",
}

var balanceSelectCmd = &cobra.Command{
	Use:   "select [ID]",
	Short: "select a balance",
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
		b, err := c.SelectBalance(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting balance")
		}
		a, err := c.SelectAccount(b.ID)
		if err != nil {
			return errors.Wrapf(err, "selecting account for balance %+v", b)
		}
		table.Accounts(storage.Accounts{*a}, os.Stdout)
		table.Balances(storage.Balances{*b}, os.Stdout)
		return nil
	},
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
		b, err := c.SelectBalance(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting balance")
		}
		a, err := c.SelectAccount(b.ID)
		if err != nil {
			return errors.Wrapf(err, "selecting account for balance %+v", b)
		}
		table.Accounts(storage.Accounts{*a}, os.Stdout)
		table.Balances(storage.Balances{*b}, os.Stdout)
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
		balanceSelectCmd,
	} {
		err := viper.BindPFlags(c.Flags())
		if err != nil {
			log.Fatal(errors.Wrap(err, "binding pflags"))
		}
		balanceCmd.AddCommand(c)
	}
}
