package statehttp

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type stateHandler struct {
	state string

	tracker *responseTracker

	writer  http.ResponseWriter
	request *http.Request
}

type responseTracker struct {
	state  bool
	metric *prometheus.CounterVec
}

func NewHandler(state string) stateHandler {
	responseStatus := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rhys_go_webserver_request_total",
		Help: "The total number of processed requests by response status by code",
	},
		[]string{"status"},
	)

	r := responseTracker{
		state:  true,
		metric: responseStatus,
	}

	return stateHandler{
		state:   state,
		tracker: &r,
	}
}

func (s stateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.writer = w
	s.request = r
	s.handleState()
}

func (s *stateHandler) handleState() {
	switch s.state {
	case "reset":
		s.writer.WriteHeader(http.StatusResetContent)
	case "ok":
		s.ok()
	case "fail":
		s.fail()
	case "both":
		if s.tracker.state == true {
			s.ok()
		}
		if s.tracker.state == false {
			s.fail()
		}
		s.tracker.state = !s.tracker.state
	}
}

func (s *stateHandler) ok() {
	s.writer.WriteHeader(http.StatusOK)
	s.writer.Write([]byte("200 - Something good happened! \n"))
	s.tracker.metric.WithLabelValues("200").Inc()
}

func (s *stateHandler) fail() {
	s.writer.WriteHeader(http.StatusInternalServerError)
	s.writer.Write([]byte("500 - Something bad happened! \n"))
	s.tracker.metric.WithLabelValues("500").Inc()
}
