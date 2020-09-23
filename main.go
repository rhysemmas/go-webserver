package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	respStatus = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rhys_go_webserver_request_total",
		Help: "The total number of processed requests by response status by code",
	},
		[]string{"status"},
	)
)

func main() {
	addr := os.Getenv("ADDR")
	state := os.Getenv("STATE")

	if state == "" {
		log.Fatalf("STATE env var must be set to 'reset', ok', 'fail', or 'both' (for switching between 200 and 500)")
	}

	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responseSuccess := false

		if state == "reset" {
			w.WriteHeader(http.StatusResetContent)
		}

		if state == "ok" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 - Something good happened v2! \n"))
			respStatus.WithLabelValues("200").Inc()
		}

		if state == "fail" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened v2! \n"))
			respStatus.WithLabelValues("500").Inc()
		}

		if state == "both" {
			responseSuccess = !responseSuccess
			if responseSuccess == true {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("200 - Something good happened v2! \n"))
				respStatus.WithLabelValues("200").Inc()
			}
			if responseSuccess == false {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Something bad happened v2! \n"))
				respStatus.WithLabelValues("500").Inc()
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
