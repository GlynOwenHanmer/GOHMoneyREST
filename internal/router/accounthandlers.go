package router

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/mon/internal/model"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// TODO: redesign these so that they don't need to take a request? There could
// TODO: be multiple handler types either take a request or don't take a request
func (env *environment) handlerSelectAccounts(_ *http.Request) (int, interface{}, error) {
	as, err := env.storage.SelectAccounts()
	if err != nil {
		return http.StatusServiceUnavailable, nil, errors.Wrap(err, "selecting Accounts from client")
	}
	return http.StatusOK, as, nil
}

func (env *environment) muxAccountIDHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}
	return env.handlerSelectAccount(id)
}

func (env *environment) handlerSelectAccount(id uint) (int, interface{}, error) {
	a, err := env.storage.SelectAccount(id)
	if err != nil {
		return http.StatusNotFound, nil, errors.Wrapf(err, "selecting Account with id:%d from storage", id)
	}
	return http.StatusOK, a, nil
}

func (env *environment) muxAccountInsertHandlerFunc(r *http.Request) (int, interface{}, error) {
	bod, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "reading request body")
	}

	defer func() {
		// TODO: this handler only needs to take a []byte which would mean we
		// TODO: can handle closing the body elsewhere
		cErr := r.Body.Close()
		if cErr != nil {
			log.Print(errors.Wrap(err, "closing request body"))
		}
	}()

	a, err := account.UnmarshalJSON(bod)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "unmarshalling request body")
	}
	return env.handlerInsertAccount(*a)
}

func (env *environment) handlerInsertAccount(a account.Account) (int, interface{}, error) {
	inserted, err := env.storage.InsertAccount(a)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "inserting Account into storage")
	}
	return http.StatusOK, inserted, nil
}

func (env *environment) muxAccountUpdateHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}

	o, err := env.storage.SelectAccount(id)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "selecting account with id:%d", id)
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

	updates, err := account.UnmarshalJSON(bod)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "unmarshalling request body")
	}

	return env.handlerUpdateAccount(*o, *updates)
}

func (env *environment) handlerUpdateAccount(a storage.Account, updates account.Account) (int, interface{}, error) {
	updated, err := model.UpdateAccount(env.storage, a, updates)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	return http.StatusOK, updated, nil
}

func (env *environment) muxAccountDeleteHandlerFunc(r *http.Request) (int, interface{}, error) {
	id, err := extractID(mux.Vars(r))
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "extracting account ID")
	}
	return env.handlerDeleteAccount(id)
}

func (env *environment) handlerDeleteAccount(id uint) (int, interface{}, error) {
	err := model.DeleteAccount(env.storage, id)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrapf(err, "deleting Account with id:%d from storage", id)
	}
	return http.StatusOK, nil, nil
}

func extractID(vars map[string]string) (uint, error) {
	if vars == nil {
		return 0, errors.New("nil vars map")
	}
	idString, ok := vars["id"]
	if !ok {
		return 0, errors.New("no account id context variable")
	}
	id, err := strconv.ParseUint(idString, 10, 64)
	return uint(id), errors.Wrapf(err, "parsing %s to uint", idString)
}
