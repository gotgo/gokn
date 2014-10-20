package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotgo/gokn/testing"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", EchoHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":23000", nil))
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	re := &testing.EchoRequest{
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
