package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
)

const (
	balancesFieldAccountID = "account_id"
	balancesFieldAmount    = "amount"
	balancesFieldID        = "id"
	balancesFieldTime      = "time"
	balancesTable          = "balances"
)

var (
	balancesSelectFields = fmt.Sprintf(
		"%s, %s, %s",
		balancesFieldID,
		balancesFieldTime,
		balancesFieldAmount)

	balancesSelectPrefix = fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s IS NULL `,
		balancesSelectFields,
		balancesTable,
		fieldDeleted)

	balancesSelectBalancesForAccountID = fmt.Sprintf(
		"%sAND %s = $1 ORDER BY %s ASC, %s ASC;",
		balancesSelectPrefix,
		balancesFieldAccountID,
		balancesFieldTime,
		balancesFieldID)

	balancesInsertFields = fmt.Sprintf(
		"%s, %s, %s",
		balancesFieldAccountID,
		balancesFieldTime,
		balancesFieldAmount)

	balancesInsertBalance = fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1, $2, $3) RETURNING %s;`,
		balancesTable,
		balancesInsertFields,
		balancesSelectFields)

	balancesDeleteBalance = fmt.Sprintf(
		`UPDATE %s SET %s = $1 WHERE id = $2;`,
		balancesTable,
		fieldDeleted,
	)
)

// SelectAccountBalances returns all Balances for a given account ID and any
// errors that occur whilst attempting to retrieve the Balances. The Balances
// are sorted by chronological order then by the id of the Balance in the DB
func (pg postgres) SelectAccountBalances(id uint) (*storage.Balances, error) {
	return queryBalances(pg.db, balancesSelectBalancesForAccountID, id)
}

// SelectBalanceByAccountAndID selects a balance with a given ID within a given
// account. An error will be returned if no balance can be found with the ID
// for the given account.
func (pg postgres) SelectBalanceByAccountAndID(a storage.Account, balanceID uint) (*storage.Balance, error) {
	bs, err := pg.SelectAccountBalances(a.ID)
	if err != nil {
		return nil, errors.Wrap(err, "selecting account balances for account %+v")
	}
	for _, b := range *bs {
		if b.ID == balanceID {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("no balance with id %d for account", balanceID)
}

func (pg postgres) InsertBalance(accountID uint, b balance.Balance, _ string) (*storage.Balance, error) {
	return queryBalance(pg.db, balancesInsertBalance, accountID, b.Date, b.Amount)
}

func (pg postgres) DeleteBalance(id uint) error {
	_, err := queryBalance(pg.db, balancesDeleteBalance, time.Now(), id)
	return errors.Wrap(err, "querying balance")
}

// queryBalance returns an error if more than one result is returned from the query
// queryBalance may or may not return an error if zero results are returned.
func queryBalance(db *sql.DB, queryString string, values ...interface{}) (*storage.Balance, error) {
	bs, err := queryBalances(db, queryString, values...)
	if err != nil {
		return nil, errors.Wrap(err, "querying balances")
	}
	if len(*bs) > 1 {
		return nil, errors.New("query returned more than 1 result")
	}
	if bs == nil || len(*bs) == 0 {
		return nil, nil
	}
	return &(*bs)[0], nil
}

func queryBalances(db *sql.DB, queryString string, values ...interface{}) (*storage.Balances, error) {
	rows, err := db.Query(queryString, values...)
	if err != nil {
		return nil, errors.Wrap(err, "querying db")
	}
	defer nonReturningCloseRows(rows)
	return scanRowsForBalances(rows)
}

// scanRowsForBalance scans a sql.Rows for a Balances object and returns any
// error occurring along the way.
func scanRowsForBalances(rows *sql.Rows) (bs *storage.Balances, err error) {
	bs = &storage.Balances{}
	for rows.Next() {
		var ID uint
		var date time.Time
		var amount float64
		err = rows.Scan(&ID, &date, &amount)
		if err != nil {
			return nil, errors.Wrap(err, "scanning rows")
		}
		var innerB *balance.Balance
		innerB, err = balance.New(date, balance.Amount(int(amount)))
		if err != nil {
			return nil, errors.Wrap(err, "creating new balance from scan results")
		}
		*bs = append(*bs, storage.Balance{ID: ID, Balance: *innerB})
	}
	if err == nil {
		err = errors.Wrap(rows.Err(), "rows error: ")
	}
	return
}
