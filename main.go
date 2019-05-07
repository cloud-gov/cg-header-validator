package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/18f/header-validator/libheader"
	"github.com/cloudfoundry-community/go-cfenv"
	"log"
	"net/http"
	"os"
)

var (
	defaultPort  = 8080
	comparer    *libheader.HeaderComparer
	filename    string
)

func main() {
	log.Print("hello from boulder.")

	flag.StringVar(&filename, "header-ref", "", "A JSON filename in the format of map[string][]string to use as the expected headers.")
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	comparer = libheader.NewComparer()
	if ok, err := comparer.Load(file); !ok {
		log.Fatal(err)
	}

	var port string
	appEnv, err := cfenv.Current()
	if err != nil {
		port = fmt.Sprintf(":%d", defaultPort)
	} else {
		port = fmt.Sprintf(":%d", appEnv.Port)
	}

	http.HandleFunc("/", headerRepeater)
	http.HandleFunc("/diff", headerDiff)
	log.Fatal(http.ListenAndServe(port, nil))
}

func headerRepeater(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(r.Header)
}

func headerDiff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(comparer.Compare(r.Header))
}
