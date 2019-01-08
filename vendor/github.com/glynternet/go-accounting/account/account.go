package account

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/go-money/currency"
	gtime "github.com/glynternet/go-time"
	"github.com/pkg/errors"
)

// New creates a new Account object with a given name, currency.Code and start
// time.
// New returns the created Account or an error if any of the Account parameters are invalid
func New(name string, currencyCode currency.Code, opened time.Time, os ...Option) (*Account, error) {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return nil, errors.New(EmptyNameError)
	}
	a := &Account{
		name:         trimmed,
		currencyCode: currencyCode,
	}
	err := gtime.Start(opened)(&a.timeRange)
	if err != nil {
		return nil, err
	}
	for _, o := range os {
		if o == nil {
			continue
		}
		err = o(a)
		if err != nil {
			return nil, err
		}
	}
	err = a.validate()
	if err != nil {
		a = nil
	}
	return a, err
}

// Account holds the logic for an financial account.
type Account struct {
	name         string
	timeRange    gtime.Range
	currencyCode currency.Code
}

// Name returns the name associated with a given Account.
func (a Account) Name() string {
	return a.name
}

// Opened returns the start time that the Account opened.
func (a Account) Opened() time.Time {
	return a.timeRange.Start().Time
}

// Closed returns the a NullTime object that is Valid if the Account has been closed.
func (a Account) Closed() gtime.NullTime {
	return a.timeRange.End()
}

// TimeRange returns the time range that represents the lifetime of the Account.
func (a Account) TimeRange() gtime.Range {
	return a.timeRange
}

// OpenAt returns true if the Account is open at a given time.
// If an account is closed at the exact time as the reference time, OpenAt will
// return false, including when an account was opened and closed at exactly the
// same time.
func (a Account) OpenAt(t time.Time) bool {
	if a.Opened().After(t) {
		return false
	}
	if a.Closed().Valid && !a.Closed().Time.After(t) {
		return false
	}
	return true
}

// CurrencyCode returns the currency code of the Account.
func (a Account) CurrencyCode() currency.Code {
	return a.currencyCode
}

// validate checks the state of an Account to see if it is has any logical errors. validate returns a set of errors representing errors with different fields of the Account.
func (a Account) validate() (err error) {
	var fieldErrorDescriptions []string
	if len(a.name) == 0 {
		fieldErrorDescriptions = append(fieldErrorDescriptions, EmptyNameError)
	}
	if len(fieldErrorDescriptions) > 0 {
		err = FieldError(fieldErrorDescriptions)
	}
	return
}

// ValidateBalance validates a balance against an Account.
// ValidateBalance returns any logical errors between the Account and the balance.
// ValidateBalance first attempts to validate the Account as an entity by itself. If there are any errors with the Account, these errors are returned and the balance is not attempted to be validated against the Account.
// If the date of the balance is outside of the TimeRange of the Account, a DateOutOfAccountTimeRange will be returned.
func (a Account) ValidateBalance(b balance.Balance) (err error) {
	err = a.validate()
	if err != nil {
		return
	}
	if !a.timeRange.Contains(b.Date) && (!a.Closed().Valid || !a.Closed().Time.Equal(b.Date)) {
		return balance.DateOutOfAccountTimeRange{
			BalanceDate:      b.Date,
			AccountTimeRange: a.timeRange,
		}
	}
	return
}

// MarshalJSON marshals an Account into a json blob, returning the blob with any errors that occur during the marshalling.
func (a Account) MarshalJSON() ([]byte, error) {
	type Alias Account
	return json.Marshal(&struct {
		*Alias
		Name     string
		Opened   time.Time
		Closed   gtime.NullTime
		Currency currency.Code
	}{
		Alias:    (*Alias)(&a),
		Name:     a.Name(),
		Opened:   a.Opened(),
		Closed:   a.Closed(),
		Currency: a.currencyCode,
	})
}

// UnmarshalJSON attempts to unmarshal a json blob into an Account object,
// returning any errors that occur during the unmarshalling.
func (a *Account) UnmarshalJSON(data []byte) (err error) {
	type Alias Account
	aux := &struct {
		Name     string
		Opened   time.Time
		Closed   gtime.NullTime
		Currency string
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "unmarshalling data to auxilliary")
	}
	a.name = aux.Name
	c, err := currency.NewCode(aux.Currency)
	if err != nil {
		return errors.Wrapf(err, "creating new currency for %s", aux.Currency)
	}
	a.currencyCode = *c
	tr := new(gtime.Range)
	err = gtime.Start(aux.Opened)(tr)
	if err != nil {
		return errors.Wrap(err, "applying Account start time")
	}
	if aux.Closed.Valid {
		err = gtime.End(aux.Closed.Time)(tr)
		if err != nil {
			return errors.Wrap(err, "applying Account end time")
		}
	}
	a.timeRange = *tr
	return a.validate()
}

// UnmarshalJSON unmarshals json data into a struct that satisfies the Account interface
func UnmarshalJSON(data []byte) (*Account, error) {
	var a Account
	err := json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Equal returns true if both accounts a and b are logically the same.
func (a Account) Equal(b Account) bool {
	switch {
	case a.name != b.Name():
		return false
	case !a.timeRange.Equal(b.TimeRange()):
		return false
	}
	return true
}
