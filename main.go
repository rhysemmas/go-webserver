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
		addr = ":8181"
	}

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/", stateHandler(state))
	mux.HandleFunc("/ok", okHandler)
	mux.HandleFunc("/fail", failHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting server on %s", addr)
	log.Println(server.ListenAndServe())
}

func stateHandler(state string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseSuccess := false

		switch state {
		case "reset":
			w.WriteHeader(http.StatusResetContent)
		case "ok":
			okHandler(w, r)
		case "fail":
			failHandler(w, r)
		case "both":
			responseSuccess = !responseSuccess
			if responseSuccess == true {
				okHandler(w, r)
			}
			if responseSuccess == false {
				failHandler(w, r)
			}
		}
	}
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Something good happened! \n"))
	respStatus.WithLabelValues("200").Inc()
}

func failHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Something bad happened! \n"))
	respStatus.WithLabelValues("500").Inc()

}
