package testing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotgo/fw/io"
)

func StartEchoHandler() *io.GracefulListener {

	r := mux.NewRouter()
	r.HandleFunc("/", EchoHandler)
	httpMux := http.NewServeMux()
	httpMux.Handle("/", r)
	server := &http.Server{
		Handler: httpMux,
	}

	var gracefulListener *io.GracefulListener
	port := "23000"
	if listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	} else if gracefulListener, err = io.MakeGraceful(listener); err != nil {
		panic(err)
	}
	go func() {
		server.Serve(gracefulListener)
	}()

	return gracefulListener
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	re := &EchoRequest{
		URL:     r.URL,
		Headers: r.Header,
		Method:  r.Method,
	}

	if bytes, err := ioutil.ReadAll(r.Body); err != nil {
		re.Error = err
	} else {
		var objmap map[string]interface{}
		err = json.Unmarshal(bytes, &objmap)
		re.Body = objmap
	}

	if b, err := json.MarshalIndent(re, "", "\t"); err != nil {
		panic(err)
	} else {
		w.Write(b)
	}
}
