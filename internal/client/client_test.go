package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glynternet/mon/pkg/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// ensure that a Client can be used as a storage.Storage
var _ storage.Storage = Client{url: ""}

func Test_getBodyFromEndpoint(t *testing.T) {
	t.Run("get error", func(t *testing.T) {
		c := Client{url: "bloopybloop"}
		bod, err := c.getBodyFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting from endpoint")
		assert.Nil(t, bod)
	})

	t.Run("unexpected status", func(t *testing.T) {
		srv := newJSONTestServer(nil, http.StatusTeapot)
		defer srv.Close()
		c := Client{url: srv.URL}
		as, err := c.getBodyFromEndpoint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "server returned unexpected code")
		assert.Nil(t, as)
	})
}

type stubMarshal struct {
	err error
}

func (m stubMarshal) MarshalJSON() ([]byte, error) {
	return nil, m.err
}

func Test_postAsJSONToEndpoint(t *testing.T) {
	t.Run("marshal error", func(t *testing.T) {
		c := Client{url: "bloopybloop"}
		obj := stubMarshal{
			err: errors.New("can't unmarshal me"),
		}
		res, err := c.postAsJSONToEndpoint("", obj)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "can't unmarshal me")
		}
		assert.Nil(t, res)
	})

	t.Run("post to endpoint error", func(t *testing.T) {
		c := Client{url: "bloopybleep"}
		res, err := c.postAsJSONToEndpoint("", nil)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "posting to endpoint")
		}
		assert.Nil(t, res)
	})
}

func newJSONTestServer(encode interface{}, code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs, err := json.Marshal(encode)
		if err != nil {
			panic(fmt.Sprintf("error marshalling to json: %v", err))
		}
		w.WriteHeader(code)
		w.Header().Set(`Content-Type`, `application/json; charset=UTF-8`)
		_, err = w.Write(bs)
		if err != nil {
			panic(fmt.Sprintf("error writing to ResponseWriter: %v", err))
		}
	}))
}
