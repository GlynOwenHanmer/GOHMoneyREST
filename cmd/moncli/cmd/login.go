package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/glynternet/mon/internal/monauth"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const keyAuthServerURL = "auth-server-url"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, args []string) error {
		tokenFilePath := viper.GetString(keyAuthTokenFile)
		if tokenFilePath == "" {
			return fmt.Errorf("%s is not set", keyAuthTokenFile)
		}

		authServerURL := viper.GetString(keyAuthServerURL)
		if authServerURL == "" {
			return fmt.Errorf("%s is not set", keyAuthServerURL)
		}

		authURL := fmt.Sprintf("%s/loginurl", authServerURL)

		resp, err := (&http.Client{
			Timeout: time.Second * 5,
		}).Get(authURL)

		if err != nil {
			return errors.Wrapf(err, "getting login URL from url: %s", authURL)
		}

		defer func() {
			cErr := resp.Body.Close()
			if err == nil {
				err = errors.Wrap(cErr, "closing response body")
				return
			}
			if cErr != nil {
				log.Printf("error closing response body: %+v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status:%d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		var bod monauth.LoginURLResponse
		err = json.NewDecoder(resp.Body).Decode(&bod)
		if err != nil {
			return errors.Wrap(err, "decoding response")
		}

		fmt.Printf("Go to the following url in your browser to authenticate:\n%s\n", bod.LoginURL)
		fmt.Println("If login is successful, paste the output here...")

		reader := bufio.NewReader(os.Stdin)
		callbackResponse, err := reader.ReadBytes('\n')
		if err != nil {
			return errors.Wrap(err, "reading input")
		}

		// remove delimiter
		callbackResponse = callbackResponse[:len(callbackResponse)-1]

		return errors.Wrap(ioutil.WriteFile(tokenFilePath, callbackResponse, 0600), "writing token to file")
	},
}

func init() {
	loginCmd.Flags().String(keyAuthServerURL, "", "auth server url")
	err := viper.BindPFlags(loginCmd.Flags())
	if err != nil {
		log.Fatal(errors.Wrap(err, "binding pflags"))
	}
	rootCmd.AddCommand(loginCmd)
}
