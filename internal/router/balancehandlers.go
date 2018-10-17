package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/internal/model"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (env *environment) balances(accountID uint) (int, interface{}, error) {
	a, err := env.storage.SelectAccount(accountID)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "selecting account with id %d", accountID)
	}
	var bs *storage.Balances
	bs, err = model.SelectAccountBalances(env.storage, *a)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "selecting balances for account %+v", *a)
	}
	return http.StatusOK, bs, nil
}

func (env *environment) muxAccountBalancesHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}
	return env.balances(id)
}

func (env *environment) insertBalance(accountID uint, b balance.Balance) (int, interface{}, error) {
	a, err := env.storage.SelectAccount(accountID)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "selecting account")
	}
	inserted, err := model.InsertBalance(env.storage, *a, b)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "inserting balance")
	}
	return http.StatusOK, inserted, nil
}

func (env *environment) muxBalanceDeleteHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting ID")
	}
	return env.deleteBalance(id)
}

func (env *environment) deleteBalance(id uint) (int, interface{}, error) {
	err := env.storage.DeleteBalance(id)
	if err != nil {
		return http.StatusBadRequest, "", errors.Wrap(err, "deleting balance")
	}
	return http.StatusOK, "", nil
}

func (env *environment) muxAccountBalanceInsertHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}

	bod, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "reading request body")
	}

	defer func() {
		// TODO: this handler only needs to take a []byte which would mean we can handle closing the body elsewhere
		cErr := r.Body.Close()
		if cErr != nil {
			log.Print(errors.Wrap(err, "closing request body"))
		}
	}()

	var b balance.Balance
	err = json.Unmarshal(bod, &b)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "unmarshalling request body")
	}
	return env.insertBalance(id, b)
}
