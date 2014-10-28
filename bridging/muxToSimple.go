package bridging

import (
	"net/http"

	"github.com/gorilla/mux"
)

type MuxToSimple struct {
	Router *mux.Router
}

func (mts *MuxToSimple) RegisterRoute(verb, path string, f func(http.ResponseWriter, *http.Request)) {
	mts.Router.HandleFunc(path, f).Methods(verb)
}

func (mts *MuxToSimple) RequestArgs(req *http.Request) map[string]string {
	return mux.Vars(req)
}
