package main

import (
	gerrors "errors"
	"log"
	"strconv"
	"strings"
	"time"

	"encoding/csv"
	"os"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/internal/client"
	"github.com/glynternet/mon/pkg/filter"
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
	currencyString = "EUR"
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
		cc, err := currency.NewCode(currencyString)
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
				Account:  a.Account,
				Balances: bs,
			})
		}

		var times []time.Time
		for i := -daysEitherSide; i <= daysEitherSide; i++ {
			times = append(times, now.Add(time.Hour*24*time.Duration(i)))
		}

		abs, err := recurringCostsAccounts(times)
		if err != nil {
			return errors.Wrap(err, "getting recurring costs accounts")
		}
		abss = append(abss, abs)

		datedBalances := [][]string{makeHeader(abss)}

		for _, t := range times {
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

func recurringCostsAccounts(times []time.Time) (AccountBalances, error) {
	rc, err := getRecurringCosts()
	if err != nil {
		return AccountBalances{}, errors.Wrap(err, "getting recurring costs")
	}
	abs, err := rc.generateAccountBalances(times)
	if err != nil {
		return AccountBalances{}, errors.Wrap(err, "generating recurring cost account")
	}
	log.Println(abs.Account.Name())
	return abs, nil
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

func makeHeader(accounts []AccountBalances) []string {
	hs := []string{"date"}
	for _, a := range accounts {
		hs = append(hs, a.Name())
	}
	hs = append(hs, "total")
	return hs
}

func makeRow(date time.Time, abss []AccountBalances) ([]string, error) {
	dateString := date.Format("20060102")
	row := []string{dateString}
	var bs balance.Balances
	for _, abs := range abss {
		b, err := abs.Balances.AtTime(date)
		switch {
		case err == nil:
			row = append(row, strconv.Itoa(b.Amount))
			bs = append(bs, b)
		case err.Error() == gerrors.New(balance.ErrNoBalances).Error():
			row = append(row, "")
		case err != nil:
			return nil, errors.Wrapf(err, "getting balance for account:%s at time:%s", abs.Account.Name(), dateString)
		}
	}
	total := bs.Sum()
	row = append(row, strconv.Itoa(total))
	return row, nil
}

type AccountBalances struct {
	account.Account
	balance.Balances
}

func getRecurringCosts() (recurringCost, error) {
	cc, err := currency.NewCode(currencyString)
	if err != nil {
		return dailyRecurringCost{}, errors.Wrap(err, "creating new currency code")
	}
	return dailyRecurringCost{
		name:   "daily spending",
		Code:   *cc,
		Amount: -6000,
	}, nil
}
