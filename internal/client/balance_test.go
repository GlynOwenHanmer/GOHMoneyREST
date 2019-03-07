package client

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetBalancesFromEndpoint(t *testing.T) {
	t.Run("get body error", func(t *testing.T) {
		c := Client{url: "bloopybloop"}
		as, err := c.getBalancesFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting from endpoint")
		assert.Nil(t, as)
	})

	t.Run("unmarshallable response", func(t *testing.T) {
		srv := newJSONTestServer(
			struct{ NonBalances string }{NonBalances: "bloop"},
			http.StatusOK,
		)
		defer srv.Close()
		c := Client{url: srv.URL}
		bs, err := c.getBalancesFromEndpoint("")
		if assert.Error(t, err) {
			assert.IsType(t, &json.UnmarshalTypeError{}, errors.Cause(err))
		}
		assert.Nil(t, bs)
	})
}

func TestGetBalanceFromEndpoint(t *testing.T) {
	t.Run("get body error", func(t *testing.T) {
		c := Client{url: "bloopybloop"}
		as, err := c.getBalanceFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting from endpoint")
		assert.Nil(t, as)
	})

	t.Run("unmarshallable response", func(t *testing.T) {
		srv := newJSONTestServer(
			[]struct{ NonBalance string }{
				{NonBalance: "bloop"}, {}, {},
			},
			http.StatusOK,
		)
		defer srv.Close()
		c := Client{url: srv.URL}
		bs, err := c.getBalanceFromEndpoint("")
		if assert.Error(t, err) {
			assert.IsType(t, &json.UnmarshalTypeError{}, errors.Cause(err))
		}
		assert.Nil(t, bs)
	})
}
