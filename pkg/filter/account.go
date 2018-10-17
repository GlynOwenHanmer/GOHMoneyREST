package filter

import (
	"time"

	"github.com/glynternet/go-money/currency"
	"github.com/glynternet/mon/pkg/storage"
)

// AccountCondition is a function that will return true if a given storage.Account
// satisfies some certain condition.
type AccountCondition func(storage.Account) bool

// Filter returns a set of storage.Accounts that have been filtered down to the
// ones that match the AccountCondition
func (ac AccountCondition) Filter(as storage.Accounts) storage.Accounts {
	var filtered storage.Accounts
	for _, a := range as {
		if ac(a) {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// AccountNot produces an AccountCondition that inverts the outcome of the given
// AccountCondition
func AccountNot(c AccountCondition) AccountCondition {
	return func(a storage.Account) bool {
		return !c(a)
	}
}

// Existed produces an AccountCondition that can be used to identify if an
// Account existed/exists/will-exist at a given time.
// Existed will identify that an Account existed if its open date matches or
// was before the given time
func Existed(t time.Time) AccountCondition {
	return func(a storage.Account) bool {
		return !a.Account.Opened().After(t)
	}
}

// OpenAt produces an AccountCondition that will identify if a storage.Account
// was/is/will-be open at a given time
func OpenAt(t time.Time) AccountCondition {
	return func(a storage.Account) bool {
		return a.Account.OpenAt(t)
	}
}

// Currency produces an AccountCondition that will identify a storage.Account if
// it has a given currency.Code
func Currency(c currency.Code) AccountCondition {
	return func(a storage.Account) bool {
		return a.Account.CurrencyCode() == c
	}
}

// ID produces an AccountCondition that will identify a storage.Account if it
// matches an ID
func ID(id uint) AccountCondition {
	return func(a storage.Account) bool {
		return a.ID == id
	}
}

// AccountConditions is a set of AccountCondition
type AccountConditions []AccountCondition

// Or identifies when an account satisfies one of more constraints of an
// AccountConditions
func (acs AccountConditions) Or(a storage.Account) bool {
	for _, ac := range acs {
		if ac(a) {
			return true
		}
	}
	return false
}

// And identifies when an account satisfies every AccountCondition of an
// AccountConditions.
// When no AccountConditions are present, the storage.Account will always match
func (acs AccountConditions) And(a storage.Account) bool {
	for _, ac := range acs {
		if !ac(a) {
			return false
		}
	}
	return true
}
