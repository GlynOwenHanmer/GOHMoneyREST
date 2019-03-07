package router

import (
	"log"
	"net/http"

	"github.com/glynternet/mon/internal/middleware"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	// EndpointAccounts is the endpoint for Accounts
	EndpointAccounts = "/accounts"
	patternAccounts  = EndpointAccounts

	// EndpointAccount is the base endpoint for single account requests
	EndpointAccount = "/account"

	// EndpointFmtAccount is the format string for use when generating a
	// specific Account endpoint
	// TODO: this is a code smell and should probably be provided as a return
	// TODO: value from a method on another Account type, something like
	// TODO: router.Account.Endpoint(), but probably actually something better
	// than that ahahahah
	EndpointFmtAccount = EndpointAccount + "/%d"
	patternAccount     = EndpointAccount + "/{id}"

	// EndpointAccountInsert is the endpoint for inserting an Account
	EndpointAccountInsert = EndpointAccount + "/insert"
	// EndpointFmtAccountUpdate is the format string for generating the
	// endpoint to use when updating a specific Account
	EndpointFmtAccountUpdate = EndpointFmtAccount + "/update"
	patternAccountUpdate     = patternAccount + "/update"

	// EndpointBalance is the base endpoint for single balance requests
	EndpointBalance = "/balance"

	// EndpointFmtBalance is the format string for generating single balance
	// request endpoints
	EndpointFmtBalance = EndpointBalance + "/%d"
	patternBalance     = EndpointBalance + "/{id}"

	// EndpointFmtAccountBalances is the format string for use when generating
	// the endpoint to get the balances for a specific Account
	EndpointFmtAccountBalances = EndpointAccount + "/%d/balances"
	patternAccountBalances     = EndpointAccount + "/{id}/balances"

	// EndpointFmtAccountBalanceInsert is the format string for use when generating
	// the endpoint insert a Balance for a specific Account
	EndpointFmtAccountBalanceInsert = EndpointAccount + "/%d/balance/insert"
	patternAccountBalanceInsert     = EndpointAccount + "/{id}/balance/insert"
)

// New creates a new mux.Router and initialises it with generateRoutes for the store
func New(store storage.Storage, log *log.Logger) (*mux.Router, error) {
	if store == nil {
		return nil, errors.New("nil store")
	}
	rs := generateRoutes(environment{storage: store})
	return newRouter(rs, log)
}

// New creates a new Router and initialises it will all of the global generateRoutes
func newRouter(rs []route, log *log.Logger) (*mux.Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range rs {
		handler := middleware.Logger(log, route.appHandler, route.name)
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(handler)
		log.Printf("Route registered: %+v", route)
	}
	return router, nil
}

type environment struct {
	storage storage.Storage
}

func generateRoutes(e environment) []route {
	return []route{
		{
			name:       "Accounts",
			pattern:    patternAccounts,
			appHandler: e.handlerSelectAccounts,
			method:     http.MethodGet,
		},
		{
			name:       "Account",
			pattern:    patternAccount,
			appHandler: e.muxAccountIDHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "AccountInsert",
			pattern:    EndpointAccountInsert,
			appHandler: e.muxAccountInsertHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "AccountUpdate",
			pattern:    patternAccountUpdate,
			appHandler: e.muxAccountUpdateHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "AccountDelete",
			pattern:    patternAccount,
			appHandler: e.muxAccountDeleteHandlerFunc,
			method:     http.MethodDelete,
		},
		{
			name:       "Balances",
			pattern:    patternAccountBalances,
			appHandler: e.muxAccountBalancesHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "BalanceInsert",
			pattern:    patternAccountBalanceInsert,
			appHandler: e.muxAccountBalanceInsertHandlerFunc,
			method:     http.MethodPost,
		},
		{
			name:       "BalanceSelect",
			pattern:    patternBalance,
			appHandler: e.muxBalanceSelectHandlerFunc,
			method:     http.MethodGet,
		},
		{
			name:       "BalanceDelete",
			pattern:    patternBalance,
			appHandler: e.muxBalanceDeleteHandlerFunc,
			method:     http.MethodDelete,
		},
	}
}
