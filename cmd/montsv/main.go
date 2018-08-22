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
	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "montsv"

	keyServerHost   = "server-host"
	keyHistoricDays = "historic-days"
	keyForecastDays = "forecast-days"
)

func main() {
	err := cmdTSV.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

const (
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

		var times []time.Time
		for i := -viper.GetInt(keyHistoricDays); i <= viper.GetInt(keyForecastDays); i++ {
			times = append(times, now.Add(time.Hour*24*time.Duration(i)))
		}
		if len(times) == 0 {
			return errors.New("date range yielded no dates")
		}

		*as = filter.AccountCondition(filter.AccountConditions{
			filter.Existed(times[len(times)-1]),
			func(a storage.Account) bool {
				closedBeforeFirstTime := a.Account.Closed().Valid && a.Account.Closed().Time.Before(times[0])
				return !closedBeforeFirstTime
			},
		}.And).Filter(*as)

		var abss []AccountBalances
		for _, a := range *as {
			sbs, err := c.SelectAccountBalances(a)
			if err != nil {
				errors.Wrap(err, "selecting balances for account")
			}
			var bs = sbs.InnerBalances()
			abss = append(abss, AccountBalances{
				Account:  a.Account,
				Balances: bs,
			})
		}

		gabss, err := generatedAccountBalances(times)
		if err != nil {
			return errors.Wrap(err, "getting recurring costs accounts")
		}
		abss = append(abss, gabss...)

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

		return errors.Wrap(w.Error(), "writing separated values")
	},
}

func generatedAccountBalances(times []time.Time) ([]AccountBalances, error) {
	var abss []AccountBalances
	for details, ag := range getAmountGenerators() {
		abs, err := generateAccountBalances(details, ag, times)
		if err != nil {
			return nil, errors.Wrap(err, "generating AccountBalances")
		}
		abss = append(abss, abs)
	}
	return abss, nil
}

func generateAccountBalances(ds accountDetails, ag amountGenerator, times []time.Time) (AccountBalances, error) {
	cc, err := currency.NewCode(ds.currencyString)
	if err != nil {
		return AccountBalances{}, errors.Wrapf(err, "creating new currency code")
	}

	a, err := account.New(ds.name, *cc, time.Time{}) // time/date of account is not used currently
	if err != nil {
		return AccountBalances{}, errors.Wrap(err, "creating new account")
	}

	var bs balance.Balances
	for _, t := range times {
		b, err := generateBalance(ag, t)
		if err != nil {
			return AccountBalances{}, errors.Wrapf(err, "generating balance for time:%s", t)
		}
		bs = append(bs, *b)
	}
	return AccountBalances{
		Account:  *a,
		Balances: bs,
	}, nil
}

func generateBalance(ag amountGenerator, at time.Time) (*balance.Balance, error) {
	b, err := balance.New(at, balance.Amount(ag.generateAmount(at)))
	return b, errors.Wrap(err, "creating balance")
}

func init() {
	cobra.OnInitialize(initConfig)
	cmdTSV.Flags().StringP(keyServerHost, "H", "", "server host")
	cmdTSV.Flags().Int(keyHistoricDays, 90, "days either side of now to provide data for")
	cmdTSV.Flags().Int(keyForecastDays, 30*6, "days in the future to provide data for")
	err := viper.BindPFlags(cmdTSV.Flags())
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

type accountDetails struct {
	name           string
	currencyString string
}

func getAmountGenerators() map[accountDetails]amountGenerator {
	return map[accountDetails]amountGenerator{
		{
			name:           "daily spending",
			currencyString: "GBP",
		}: dailyRecurringAmount{
			Amount: -1500,
			from:   now,
		},
		{
			name:           "storage",
			currencyString: "EUR",
		}: monthlyRecurringCost{
			amount:      -7900,
			dateOfMonth: 1,
			from:        time.Date(2018, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "health insurance",
			currencyString: "EUR",
		}: monthlyRecurringCost{
			amount:      -10250,
			dateOfMonth: 27,
			from:        now,
		},
		{
			name:           "energy bill",
			currencyString: "EUR",
		}: monthlyRecurringCost{
			amount:      -3150,
			dateOfMonth: 12,
			from:        now,
		},
		{
			name:           "haircut",
			currencyString: "GBP",
		}: monthlyRecurringCost{
			amount:      -2400, // ~ every 6 weeks
			dateOfMonth: 26,
			from:        now,
		},
		{
			name:           "ABN Amro bank account",
			currencyString: "EUR",
		}: monthlyRecurringCost{
			amount:      -155, //every 6 weeks
			dateOfMonth: 19,
			from:        now,
		},
		{
			name:           "ABN Maandpremie",
			currencyString: "EUR",
		}: monthlyRecurringCost{
			amount:      -1461,
			dateOfMonth: 3,
			from:        now,
		},
		{
			name:           "O2 Phone Bill",
			currencyString: "GBP",
		}: monthlyRecurringCost{
			amount:      -3000,
			dateOfMonth: 17,
			from:        now,
		},
		{
			name:           "John & Emily Registration",
			currencyString: "EUR",
		}: dailyRecurringAmount{
			Amount: -142, // =-1000/7, Tenner a week
			from:   now,
		},
	}
}
