package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type printer struct{}

func (p printer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	defer w.WriteHeader(http.StatusOK)
	r.ParseForm()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	var o map[string]interface{}
	if len(b) > 0 {
		if err := json.Unmarshal(b, &o); err != nil {
			fmt.Printf("body not json: (%s)\n", err)
		}
	}
	b, _ = json.MarshalIndent(o, "", "  ")
	fmt.Println("-----request-----")
	fmt.Printf("%s %s %s\n", r.Proto, r.Method, r.URL.String())
	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ";"))
	}
	fmt.Println("------form-------")
	fmt.Printf("%+v\n", r.Form)
	fmt.Println("------body-------")
	fmt.Printf("%s\n", string(b))
}
