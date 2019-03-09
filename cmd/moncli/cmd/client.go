package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/glynternet/mon/internal/client"
	"github.com/glynternet/mon/internal/monauth"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewClient() (client.Client, error) {
	var options []client.Option
	tokenPath := viper.GetString(keyAuthTokenFile)
	if tokenPath != "" {
		token, err := loadToken(tokenPath)
		if err != nil {
			return client.Client{},
				errors.WithMessage(err, "error loading token")
		}
		options = append(options, client.WithAuthToken(token))
	}
	return client.New(viper.GetString(keyServerHost), options...), nil
}

func loadToken(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.WithMessage(
			err, fmt.Sprintf("reading file at path: %s", path),
		)
	}

	var token monauth.LoginCallbackResponse

	d := json.NewDecoder(bytes.NewReader(bs))
	d.DisallowUnknownFields()
	err = d.Decode(&token)
	if err != nil {
		return "", errors.WithMessage(err, "unmarshalling token")
	}

	return token.Token, nil
}
