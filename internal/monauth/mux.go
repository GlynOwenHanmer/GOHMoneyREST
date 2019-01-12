package monauth

import (
	"net/http"
	"time"

	"github.com/glynternet/mon/internal/router"
)

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
