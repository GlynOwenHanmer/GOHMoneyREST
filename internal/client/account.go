package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/glynternet/go-accounting/account"
	"github.com/glynternet/mon/internal/router"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
)

// SelectAccounts is used to retrieve accounts from the mon server
func (c Client) SelectAccounts() (*storage.Accounts, error) {
	return c.getAccountsFromEndpoint(router.EndpointAccounts)
}

func (c Client) getAccountsFromEndpoint(e string) (*storage.Accounts, error) {
	bod, err := c.getBodyFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting body from endpoint")
	}
	as := &storage.Accounts{}
	err = errors.Wrapf(json.Unmarshal(bod, as), "unmarshalling response body: %s", string(bod))
	if err != nil {
		as = nil
	}
	return as, err
}

// SelectAccount retrieves an account from the mon server by a given ID
func (c Client) SelectAccount(ID uint) (*storage.Account, error) {
	return c.getAccountFromEndpoint(fmt.Sprintf(router.EndpointFmtAccount, ID))
}

func (c Client) getAccountFromEndpoint(e string) (*storage.Account, error) {
	bod, err := c.getBodyFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting body from endpoint")
	}
	return unmarshalJSONToAccount(bod)
}

// InsertAccount will attempt to insert an account by calling the mon server and return the stored Account
func (c Client) InsertAccount(a account.Account) (*storage.Account, error) {
	bs, err := c.postAccountToEndpoint(router.EndpointAccountInsert, a)
	if err != nil {
		return nil, errors.Wrapf(err, "posting account to endpoint %s", router.EndpointAccountInsert)
	}
	return unmarshalJSONToAccount(bs)
}

// UpdateAccount will updated a currently stored account with updates provided by another account
func (c Client) UpdateAccount(id uint, updates account.Account) (*storage.Account, error) {
	endpoint := fmt.Sprintf(router.EndpointFmtAccountUpdate, id)
	bs, err := c.postAccountToEndpoint(endpoint, updates)
	if err != nil {
		return nil, errors.Wrapf(err, "posting account to endpoint %s", endpoint)
	}
	return unmarshalJSONToAccount(bs)
}

// DeleteAccount will attempt to delete an account through the mon server by the given id
func (c Client) DeleteAccount(id uint) error {
	endpoint := fmt.Sprintf(router.EndpointFmtAccount, id)
	r, err := c.deleteToEndpoint(endpoint)
	if err != nil {
		return errors.Wrapf(err, "deleting account to endpoint %s", endpoint)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d (%s)", r.StatusCode, http.StatusText(r.StatusCode))
	}
	return nil
}

func (c Client) postAccountToEndpoint(e string, a account.Account) ([]byte, error) {
	res, err := c.postAsJSONToEndpoint(e, a)
	if err != nil {
		return nil, errors.Wrap(err, "posting as JSON")
	}
	return processResponseForBody(res)
}

func unmarshalJSONToAccount(data []byte) (*storage.Account, error) {
	if len(data) == 0 {
		return nil, errors.New("no data provided")
	}
	if string(data) == "null" {
		return nil, errors.New(`data was "null"`)
	}
	a := &storage.Account{}
	err := errors.Wrapf(json.Unmarshal(data, a), "json unmarshalling into account. bytes as string: %s", data)
	if err != nil {
		a = nil
	}
	return a, err
}
