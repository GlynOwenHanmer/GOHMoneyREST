package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/glynternet/go-accounting-storage"
	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/go-money/currency"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	fieldID       = "id"
	fieldName     = "name"
	fieldOpened   = "opened"
	fieldClosed   = "closed"
	fieldCurrency = "currency"
	fieldDeleted  = "deleted"
	table         = "accounts"
)

var (
	fieldsInsert = fmt.Sprintf(
		"%s, %s, %s, %s",
		fieldName,
		fieldOpened,
		fieldClosed,
		fieldCurrency)

	fieldsSelect = fmt.Sprintf(
		"%s, %s, %s, %s, %s, %s",
		fieldID,
		fieldName,
		fieldOpened,
		fieldClosed,
		fieldCurrency,
		fieldDeleted)

	querySelectAccounts = fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s IS NULL ORDER BY %s ASC;",
		fieldsSelect,
		table,
		fieldDeleted,
		fieldID)

	querySelectAccount = fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = $1 AND %s IS NULL;",
		fieldsSelect,
		table,
		fieldID,
		fieldDeleted)

	queryInsertAccount = fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1, $2, $3, $4) returning %s`,
		table,
		fieldsInsert,
		fieldsSelect)

	queryUpdateAccount = fmt.Sprintf(
		`UPDATE %s SET %s = $1, %s = $2, %s = $3, %s = $4 WHERE %s = $5 returning %s`,
		table,
		fieldName,
		fieldOpened,
		fieldClosed,
		fieldCurrency,
		fieldID,
		fieldsSelect)
)

// SelectAccounts returns an Accounts item holding all Account entries within
// the given database along with any errors that occurred whilst attempting to
// retrieve the Accounts.
func (pg postgres) SelectAccounts() (*storage.Accounts, error) {
	return queryAccounts(pg.db, querySelectAccounts)
}

// SelectAccount returns an Account with the given id.
func (pg postgres) SelectAccount(id uint) (*storage.Account, error) {
	dba, err := queryAccount(pg.db, querySelectAccount, id)
	return dba, errors.Wrap(err, "querying Account")
}

// InsertAccount inserts an account.Account in the storage backend and returns it.
func (pg postgres) InsertAccount(a account.Account) (*storage.Account, error) {
	dba, err := queryAccount(pg.db, queryInsertAccount, a.Name(), a.Opened(), pq.NullTime(a.Closed()), a.CurrencyCode())
	return dba, errors.Wrap(err, "querying Account")
}

// UpdateAccount updates a stored account to reflect the details of some other
// account data. The updates will be verified to ensure that any data to be
// used will be logically sound with the balances and other account details.
func (pg postgres) UpdateAccount(a *storage.Account, updates *account.Account) (*storage.Account, error) {
	bs, err := pg.SelectAccountBalances(*a)
	if err != nil {
		return nil, errors.Wrap(err, "selecting Account Balances for update validation")
	}
	for _, b := range *bs {
		err := updates.ValidateBalance(b.Balance)
		if err != nil {
			return nil, fmt.Errorf("update would make balance invalid: %v", err)
		}
	}
	dba, err := queryAccount(
		pg.db,
		queryUpdateAccount,
		updates.Name(),
		updates.Opened(),
		pq.NullTime(updates.Closed()),
		updates.CurrencyCode(),
		a.ID,
	)
	return dba, errors.Wrap(err, "querying Account")
}

func queryAccount(db *sql.DB, queryString string, values ...interface{}) (*storage.Account, error) {
	as, err := queryAccounts(db, queryString, values...)
	if err != nil {
		return nil, errors.Wrap(err, "querying accounts")
	}
	resLen := len(*as)
	if resLen == 0 {
		return nil, errors.New("query returned no accounts")
	}
	if resLen > 1 {
		return nil, fmt.Errorf("expected 1 account but query returned %d", resLen)
	}
	return &(*as)[0], nil
}

func queryAccounts(db *sql.DB, queryString string, values ...interface{}) (*storage.Accounts, error) {
	rows, err := db.Query(queryString, values...)
	if err != nil {
		return nil, err
	}
	defer nonReturningCloseRows(rows)
	return scanRowsForAccounts(rows)
}

// scanRowsForAccounts scans an sql.Rows object for storage.Accounts objects
// and returns then along with any error that occurs whilst attempting to scan.
func scanRowsForAccounts(rows *sql.Rows) (*storage.Accounts, error) {
	var openAccounts storage.Accounts
	for rows.Next() {
		var id uint
		var name, code string
		var opened time.Time
		var closed, deleted pq.NullTime
		// 	fieldID, fieldName, fieldOpened, fieldClosed, fieldCurrency, fieldDeleted)
		err := rows.Scan(&id, &name, &opened, &closed, &code, &deleted)
		if err != nil {
			return nil, errors.Wrap(err, "scanning row")
		}
		c, err := currency.NewCode(code)
		if err != nil {
			return nil, errors.Wrap(err, "generating new currency code")
		}
		innerAccount, err := account.New(name, *c, opened)
		if err != nil {
			return nil, errors.Wrap(err, "creating new inner account")
		}
		if closed.Valid {
			err = account.CloseTime(closed.Time)(innerAccount)
			if err != nil {
				return nil, errors.Wrap(err, "applying closed time to inner account")
			}
		}
		a := &storage.Account{ID: id, Account: *innerAccount}
		if deleted.Valid {
			err := storage.DeletedAt(deleted.Time)(a)
			if err != nil {
				return nil, errors.Wrap(err, "applying deleted time to inner account")
			}
		}
		openAccounts = append(openAccounts, *a)
	}
	return &openAccounts, rows.Err()
}
