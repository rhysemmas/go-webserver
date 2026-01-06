package apkovlhttp

import (
	"log"
	"net/http"
)

type apkovlHandler struct {
	hostToApkovl map[string]string
}

func NewHandler(hostmap map[string]string) *apkovlHandler {
	return &apkovlHandler{hostToApkovl: hostmap}
}

func (a *apkovlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got request from host: %v", r.RemoteAddr)
	//if filePath, ok := a.hostmap[r.RemoteAddr]
}
