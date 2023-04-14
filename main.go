package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rhysemmas/go-webserver/pkg/statehttp"
)

func main() {
	addr, state, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	serve(addr, state)
}

func setup() (string, string, error) {
	state := os.Getenv("STATE")
	addr := os.Getenv("ADDR")

	if state == "" {
		return "", "", fmt.Errorf("STATE env var must be set to 'reset', ok', 'fail', or 'both' (for switching between 200 and 500)")
	}

	if addr == "" {
		addr = ":8181"
	}

	return addr, state, nil
}

func serve(addr, state string) {
	mux := http.NewServeMux()

	stateHandler := statehttp.NewHandler(state)

	mux.Handle("/", stateHandler)
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting server on %s", addr)
	log.Println(server.ListenAndServe())
}
