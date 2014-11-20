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
	return start(23000, "/", EchoHandler)
}

func StartGeneratedHandler(port int, path string, dataGen func() (int, []byte)) *io.GracefulListener {
	return start(port, path, func(w http.ResponseWriter, r *http.Request) {
		statusCode, bts := dataGen()
		w.WriteHeader(statusCode)
		w.Write(bts)
	})
}

func start(port int, path string, handler func(http.ResponseWriter, *http.Request)) *io.GracefulListener {
	r := mux.NewRouter()
	r.HandleFunc(path, handler)
	httpMux := http.NewServeMux()
	httpMux.Handle("/", r)
	server := &http.Server{
		Handler: httpMux,
	}

	var gracefulListener *io.GracefulListener
	if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
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
