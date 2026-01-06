package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rhysemmas/go-webserver/pkg/apkovlhttp"
)

func main() {
	addr, ipToApkovl, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	serve(addr, ipToApkovl)
}

func setup() (string, map[string]string, error) {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	configFile := flag.String("config", "config.json", "config file containing map of ip addresses to apkovl files")
	flag.Parse()

	config, err := os.ReadFile(*configFile)
	if err != nil {
		return "", nil, fmt.Errorf("error reading config file: %v", err)
	}
	
	var ipToApkovl map[string]string
	if err := json.Unmarshal(config, &ipToApkovl); err != nil {
		return "", nil, fmt.Errorf("error unmarshalling config file: %v", err)
	}

	return addr, ipToApkovl, nil
}

func serve(addr string, ipToApkovl map[string]string) {
	mux := http.NewServeMux()
	mux.Handle("/", apkovlhttp.NewHandler(ipToApkovl))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("starting server on %s", addr)
	log.Println(server.ListenAndServe())
}
