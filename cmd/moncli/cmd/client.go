package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/glynternet/mon/internal/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func newClient() (client.Client, error) {
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
	token := string(bs)
	if token == "" {
		return "", errors.New("file is empty")
	}
	return token, nil
}
