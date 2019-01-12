package monauth

import (
	"net/http"
	"time"

	"github.com/glynternet/mon/internal/router"
)

// ServeMux creates a http.ServeMux configured with two routes:
//
// 		/loginurl 		- To retrieve a login URL, which the user should browse
// 						  to to login
//		/logincallback	- The address that the oauth2 flow should be configured
//						  to redirect to. This is the handler that exhanges the
//						  oauth2 code for an id_token
func ServeMux(ace AuthCodeExchanger, tokenExchangeTimout time.Duration) *http.ServeMux {
	mux := &http.ServeMux{}
	env := &env{
		AuthCodeExchanger:   ace,
		generateStateValue:  generateStateValue,
		tokenExchangeTimout: tokenExchangeTimout,
	}
	mux.Handle("/loginurl", router.AppJSONHandler(env.loginURLHandler))
	mux.Handle("/logincallback", router.AppJSONHandler(env.loginCallbackHandler))
	return mux
}
