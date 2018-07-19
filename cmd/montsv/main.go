package main

import (
	gerrors "errors"
	"log"
	"strconv"
	"strings"
	"time"

	"encoding/csv"
	"os"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/internal/client"
	"github.com/glynternet/mon/pkg/filter"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "montsv"

	keyServerHost = "server-host"
)

func main() {
	err := cmdTSV.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

const (
	daysEitherSide = 90
)

var now = time.Now()

var cmdTSV = &cobra.Command{
	Use:  appName,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.Client(viper.GetString(keyServerHost))
		as, err := c.SelectAccounts()
		if err != nil {
			return errors.Wrap(err, "selecting accounts")
		}
		cc, err := currency.NewCode("EUR")
		if err != nil {
			return errors.Wrap(err, "creating currency code")
		}

		*as = filter.AccountCondition(filter.AccountConditions{
			filter.Currency(*cc),
		}.And).Filter(*as)

		var abss []AccountBalances
		for _, a := range *as {
			sbs, err := c.SelectAccountBalances(a)
			if err != nil {
				errors.Wrap(err, "selecting balances for account")
			}
			var bs balance.Balances
			for _, sb := range *sbs {
				bs = append(bs, sb.Balance)
			}
			abss = append(abss, AccountBalances{
				Account:  a,
				Balances: bs,
			})
		}

		datedBalances := [][]string{makeHeader(*as)}
		for i := -daysEitherSide; i <= daysEitherSide; i++ {
			t := now.Add(time.Hour * 24 * time.Duration(i))
			row, err := makeRow(t, abss)
			if err != nil {
				return errors.Wrapf(err, "making row at time:%s", t.Format("20060102"))
			}
			datedBalances = append(datedBalances, row)
		}

		w := csv.NewWriter(os.Stdout)
		w.Comma = '\t'
		w.WriteAll(datedBalances) // calls Flush internally

		if err := w.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}

		return nil
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	cmdTSV.PersistentFlags().StringP(keyServerHost, "H", "", "server host")
	err := viper.BindPFlags(cmdTSV.PersistentFlags())
	if err != nil {
		log.Fatal(errors.Wrap(err, "binding root command flags"))
	}
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}

func makeHeader(accounts storage.Accounts) []string {
	hs := []string{"date"}
	for _, a := range accounts {
		hs = append(hs, a.Account.Name())
	}
	return hs
}

func makeRow(date time.Time, abss []AccountBalances) ([]string, error) {
	dateString := date.Format("20060102")
	row := []string{dateString}
	for _, abs := range abss {
		b, err := abs.Balances.AtTime(date)
		if err != nil && err.Error() != gerrors.New(balance.ErrNoBalances).Error() {
			return nil, errors.Wrapf(err, "getting balance for account:%s at time:%s", abs.Account.Account.Name(), dateString)
		}
		row = append(row, strconv.Itoa(b.Amount))
	}
	return row, nil
}

type AccountBalances struct {
	storage.Account
	balance.Balances
}
