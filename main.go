package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rhysemmas/go-webserver/pkg/apkovlhttp"
)

func main() {
	addr, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	serve(addr)
}

func setup() (string, error) {
	addr := os.Getenv("ADDR")

	if addr == "" {
		addr = ":8080"
	}

	return addr, nil
}

func serve(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/", apkovlhttp.NewHandler(nil))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting server on %s", addr)
	log.Println(server.ListenAndServe())
}
