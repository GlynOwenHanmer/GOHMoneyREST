package table

import (
	"fmt"
	"io"
	"strconv"

	"github.com/glynternet/go-time"
	"github.com/glynternet/mon/internal/accountbalance"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/olekukonko/tablewriter"
)

const dateFormat = `02-01-2006`

// Accounts writes a table for a set of Accounts to a given io.Writer
func Accounts(as storage.Accounts, w io.Writer) {
	t := newDefaultTable(w)
	t.SetHeader([]string{"ID", "Name", "Opened", "Closed", "Currency"})

	for _, a := range as {
		t.Append([]string{
			strconv.FormatUint(uint64(a.ID), 10),
			a.Account.Name(),
			a.Account.Opened().Format(dateFormat),
			closedString(a.Account.Closed()),
			a.Account.CurrencyCode().String(),
		})
	}
	t.Render() // Send output
}

// AccountsWithBalance writes a table for a set of Accounts with corresponding
// Balances to a given io.Writer
func AccountsWithBalance(abs []accountbalance.AccountBalance, w io.Writer) {
	t := newDefaultTable(w)
	t.SetHeader([]string{
		"ID", "Name", "Opened", "Closed", "Currency", "Balance Date", "Balance Amount",
	})

	for _, ab := range abs {
		t.Append([]string{
			strconv.FormatUint(uint64(ab.Account.ID), 10),
			ab.Account.Account.Name(),
			ab.Account.Account.Opened().Format(dateFormat),
			closedString(ab.Account.Account.Closed()),
			ab.Account.Account.CurrencyCode().String(),
			ab.Date.Format(dateFormat),
			strconv.Itoa(ab.Amount),
		})
	}
	t.Render() // Send output
}

// Balances writes a table for a given set of storage.Balances to a given io.Writer
func Balances(bs storage.Balances, w io.Writer) {
	t := newDefaultTable(w)
	t.SetHeader([]string{"ID", "Amount", "Date", "Note"})

	for _, b := range bs {
		t.Append([]string{
			strconv.FormatUint(uint64(b.ID), 10),
			strconv.Itoa(b.Amount),
			b.Date.Format(dateFormat),
			b.Note,
		})
	}

	t.Render()
}

// Basic writes grid of string data to a given io.Writer
func Basic(data [][]string, w io.Writer) error {
	if len(data) < 2 {
		return fmt.Errorf("requires at least 2 rows of data")
	}
	t := newDefaultTable(w)
	t.SetHeader(data[0])
	for _, row := range data[1:] {
		t.Append(row)
	}
	t.Render()
	return nil
}

func newDefaultTable(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
	return table
}

func closedString(t time.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(dateFormat)
}
