package apkovlhttp

import (
	"os"
	"log"
	"net/http"
	"strings"
)

type apkovlHandler struct {
	ipToApkovl map[string]string
}

func NewHandler(hostmap map[string]string) *apkovlHandler {
	testMap := make(map[string]string)
	testMap["192.168.1.1"] = "/srv/http/test/test.apkovl.tar.gz"
	return &apkovlHandler{ipToApkovl: testMap}
	//return &apkovlHandler{ipToApkovl: hostmap}
}

func (a *apkovlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("got request from host: %v", r.RemoteAddr)
	host := strings.Split(r.RemoteAddr, ":")
	ip := host[0]
	
	filePath, ok := a.ipToApkovl[ip]
	if !ok {
		log.Printf("file not found for ip: %v", ip)
		a.errorResponse(w)
		return
	}
	
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("error opening file: %v: %v", filePath, err)
		a.errorResponse(w)
		return
	}

	stat, err := file.Stat()
	if err != nil {
		log.Printf("error statting file: %v, %v", filePath, err)
	}

	b := make([]byte, stat.Size())
	i, err := file.Read(b)
	if err != nil {
		log.Printf("error reading file: %v: %v", filePath, err)
		a.errorResponse(w)
		return			
	}

	log.Printf("read %v bytes from file: %v", i, filePath)
	w.WriteHeader(200)
	i, err = w.Write(b)
	if err != nil {
		log.Printf("error writing response: %v", err)
		return
	}

	log.Printf("wrote %v bytes", i)
}

func (a *apkovlHandler) errorResponse(w http.ResponseWriter) {
	b := make([]byte, 0)
	w.WriteHeader(500)
	w.Write(b)
}
