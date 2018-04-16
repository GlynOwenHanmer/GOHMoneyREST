package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/glynternet/accounting-rest/client"
	"github.com/glynternet/accounting-rest/pkg/table"
	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	gtime "github.com/glynternet/go-time"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyDate     = "date"
	keyAmount   = "amount"
	keyName     = "name"
	keyCurrency = "currency"
	keyOpened   = "opened"
	keyClosed   = "closed"
	keyLimit    = "limit"
)

var accountCmd = &cobra.Command{
	Use: "account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no account id given")
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parsing account id")
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting account")
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)
		return nil
	},
}

var accountAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add an account",
	RunE: func(cmd *cobra.Command, args []string) error {
		cc, err := currency.NewCode(viper.GetString(keyCurrency))
		if err != nil {
			return errors.Wrap(err, "creating new currency code")
		}
		opened, err := parseNullTime(viper.GetString(keyOpened))
		if err != nil {
			return errors.Wrap(err, "parsing opened date")
		}
		if !opened.Valid {
			opened = gtime.NullTime{
				Valid: true,
				Time:  time.Now(),
			}
		}

		closed, err := parseNullTime(viper.GetString(keyClosed))
		if err != nil {
			return errors.Wrap(err, "parsing closed date")
		}

		var ops []account.Option
		if closed.Valid {
			ops = append(ops, account.CloseTime(closed.Time))
		}
		a, err := account.New(
			viper.GetString(keyName),
			*cc,
			opened.Time,
			ops...,
		)
		if err != nil {
			return errors.Wrap(err, "creating new account for insert")
		}

		i, err := client.Client(viper.GetString(keyServerHost)).InsertAccount(a)
		if err != nil {
			return errors.Wrap(err, "inserting new account")
		}
		table.Accounts(storage.Accounts{*i}, os.Stdout)
		return nil
	},
}

var accountBalancesCmd = &cobra.Command{
	Use: "balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no account id given")
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parsing account id")
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting account")
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)

		bs, err := c.SelectAccountBalances(*a)
		if err != nil {
			return errors.Wrap(err, "selecting account balances")
		}

		limit := viper.GetInt(keyLimit)
		if limit > len(*bs) {
			limit = len(*bs)
		}
		if limit != 0 {
			*bs = (*bs)[len(*bs)-limit:]
		}

		table.Balances(*bs, os.Stdout)
		return nil
	},
}

var accountBalanceInsertCmd = &cobra.Command{
	Use: "balance-insert",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument, receieved %d", len(args))
		}
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "parsing account id")
		}
		c := client.Client(viper.GetString(keyServerHost))
		a, err := c.SelectAccount(uint(id))
		if err != nil {
			return errors.Wrap(err, "selecting account")
		}

		nt, err := parseNullTime(viper.GetString(keyDate))
		if err != nil {
			return errors.Wrap(err, "getting date")
		}

		t := time.Now()
		if nt.Valid {
			t = nt.Time
		}

		b, err := c.InsertBalance(*a, balance.Balance{
			Date:   t,
			Amount: viper.GetInt(keyAmount),
		})
		if err != nil {
			return errors.Wrap(err, "inserting balance")
		}

		table.Accounts(storage.Accounts{*a}, os.Stdout)
		table.Balances(storage.Balances{*b}, os.Stdout)
		return nil
	},
}

func parseNullTime(ds string) (gtime.NullTime, error) {
	ds = strings.TrimSpace(ds)
	if ds == "" {
		return gtime.NullTime{}, nil
	}
	t, err := parseDateString(ds)
	if err != nil {
		return gtime.NullTime{}, errors.Wrap(err, "parsing date string")
	}
	return gtime.NullTime{
		Valid: true,
		Time:  t,
	}, nil
}

func parseDateString(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02", dateString)
}

func init() {
	accountAddCmd.Flags().StringP(keyName, "n", "", "")
	accountAddCmd.Flags().StringP(keyOpened, "o", "", "opened date")
	accountAddCmd.Flags().StringP(keyClosed, "c", "", "closed date")
	accountAddCmd.Flags().String(keyCurrency, "EUR", "")

	accountBalancesCmd.Flags().UintP(keyLimit, "l", 0, "limit results")

	accountBalanceInsertCmd.Flags().StringP(keyDate, "d", "", "date of balance to insert")
	accountBalanceInsertCmd.Flags().IntP(keyAmount, "a", 0, "amount of balance to insert")

	rootCmd.AddCommand(accountCmd)
	for _, c := range []*cobra.Command{
		accountAddCmd,
		accountBalancesCmd,
		accountBalanceInsertCmd,
	} {
		err := viper.BindPFlags(c.Flags())
		if err != nil {
			log.Fatal(errors.Wrap(err, "binding pflags"))
		}
		accountCmd.AddCommand(c)
	}
}
