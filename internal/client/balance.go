package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/glynternet/go-accounting/balance"
	"github.com/glynternet/mon/internal/router"
	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
)

// SelectAccountBalances will select the Balances that are stored for a given Account
func (c Client) SelectAccountBalances(id uint) (*storage.Balances, error) {
	return c.getBalancesFromEndpoint(fmt.Sprintf(router.EndpointFmtAccountBalances, id))
}

func (c Client) getBalancesFromEndpoint(e string) (*storage.Balances, error) {
	bod, err := c.getBodyFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting body from endpoint")
	}
	bs := &storage.Balances{}
	err = errors.Wrap(json.Unmarshal(bod, bs), "unmarshalling response")
	if err != nil {
		bs = nil
	}
	return bs, err
}

// InsertBalance will insert a balance for a given Account
func (c Client) InsertBalance(accountID uint, b balance.Balance, note string) (*storage.Balance, error) {
	endpoint := fmt.Sprintf(router.EndpointFmtAccountBalanceInsert, accountID)

	res, err := c.postAsJSONToEndpoint(endpoint, router.BalanceInsertBody{
		Balance: b,
		Note:    note,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "posting BalanceInsertBody to endpoint:%s", endpoint)
	}
	bs, err := processResponseForBody(res)
	if err != nil {
		return nil, errors.Wrap(err, "processing response for body")
	}
	return unmarshalJSONToBalance(bs)
}

// SelectBalance selects a balance with a given id
func (c Client) SelectBalance(id uint) (*storage.Balance, error) {
	return nil, errors.New("not implemented")
}

// DeleteBalance deletes a balance at a given id
func (c Client) DeleteBalance(id uint) error {
	endpoint := fmt.Sprintf(router.EndpointFmtBalance, id)
	r, err := c.deleteToEndpoint(endpoint)
	if err != nil {
		return errors.Wrapf(err, "deleting balance to endpoint %s", endpoint)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d (%s)", r.StatusCode, http.StatusText(r.StatusCode))
	}
	return nil
}

func unmarshalJSONToBalance(data []byte) (*storage.Balance, error) {
	b := &storage.Balance{}
	err := errors.Wrapf(json.Unmarshal(data, b), "json unmarshalling into balance. bytes as string: %s", data)
	if err != nil {
		b = nil
	}
	return b, err
}
