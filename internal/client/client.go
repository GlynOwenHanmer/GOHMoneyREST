package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// New creates a new Client configured to communicate with the server at the given url.
func New(url string) Client {
	return Client{url: url}
}

// Client communicates with the server
type Client struct {
	url string
}

// newClient provides the Client that should be used to make any calls against
// the mon server
func newClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

func (c Client) getFromEndpoint(endpoint string) (*http.Response, error) {
	url := c.url + endpoint
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "creating new request for url:%q", url)
	}
	return newClient().Do(r)
}

func (c Client) postToEndpoint(endpoint string, contentType string, body io.Reader) (*http.Response, error) {
	url := c.url + endpoint
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "creating new request for url:%q", url)
	}
	r.Header.Set("Content-Type", contentType)
	return newClient().Do(r)
}

func (c Client) deleteToEndpoint(endpoint string) (*http.Response, error) {
	url := c.url + endpoint
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "creating new request for url:%q", url)
	}
	return newClient().Do(r)
}

// Available reports whether the mon server is available using the Client
func (c Client) Available() bool {
	// TODO: Deprecate Available in favour of something that returns more information
	_, err := c.SelectAccounts()
	return err == nil
}

// Close is a noop closer as there is not behaviour required to close this Client
func (c Client) Close() error {
	return nil
}

func (c Client) getBodyFromEndpoint(e string) ([]byte, error) {
	res, err := c.getFromEndpoint(e)
	if err != nil {
		return nil, errors.Wrap(err, "getting from endpoint")
	}
	return processResponseForBody(res)
}

func processResponseForBody(r *http.Response) ([]byte, error) {
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned unexpected code %d (%s)", r.StatusCode, r.Status)
	}
	bod, err := ioutil.ReadAll(r.Body)

	defer func() {
		// TODO: this handler only needs to take a []byte which would mean we can handle closing the body elsewhere
		cErr := r.Body.Close()
		if cErr != nil {
			log.Print(errors.Wrap(err, "closing response body"))
		}
	}()

	return bod, errors.Wrap(err, "reading response body")
}

func (c Client) postAsJSONToEndpoint(e string, thing interface{}) (*http.Response, error) {
	bs, err := json.Marshal(thing)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling json")
	}
	res, err := c.postToEndpoint(e, `application/json; charset=UTF-8`, bytes.NewReader(bs))
	return res, errors.Wrap(err, "posting to endpoint")
}
