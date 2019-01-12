package monauth

import (
	"net/http"

	"github.com/glynternet/mon/internal/router"
)

func ServeMux(ace AuthCodeExchanger) *http.ServeMux {
	mux := &http.ServeMux{}
	env := &env{
		AuthCodeExchanger:  ace,
		generateStateValue: generateStateValue,
	}
	mux.Handle("/loginurl", router.AppJSONHandler(env.loginURLHandler))
	mux.Handle("/logincallback", router.AppJSONHandler(env.loginCallbackHandler))
	return mux
}
