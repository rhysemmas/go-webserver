package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("ADDR")
	state := os.Getenv("STATE")

	if state != "ok" || state != "fail" || state != "both" || state != "reset" {
		log.Fatalf("STATE env var must be set to 'ok', 'fail', or 'both'")
	}

	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()

	responseSuccess := false

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if state == "reset" {
			w.WriteHeader(http.StatusResetContent)
		}

		if state == "ok" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 - Something good happened! \n"))
		}

		if state == "fail" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened! \n"))
		}

		if state == "both" {
			responseSuccess = !responseSuccess
			if responseSuccess == true {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("200 - Something good happened! \n"))
			}
			if responseSuccess == false {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened! \n"))
			}
		}
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting server on %s", addr)
	log.Println(server.ListenAndServe())
}
