/*

Wrapper around http.ResponseWriter to capture egress traffic

Inspired from Negroni : https://github.com/codegangsta/negroni/blob/master/response_writer.go
 */
package main

import (
	"log"
	"time"

	"net/http"
	"net/http/httputil"

	storage "github.com/ObjectIsAdvantag/smartproxy/storage"
)


var db *storage.TrafficStorage = storage.VolatileTrafficStorage()


func CreateTrafficDumper(proxy *httputil.ReverseProxy) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestBytes, err := httputil.DumpRequest(r, true)
		if err == nil {
			trace := db.CreateTrace()
			trace.Start = start
			trace.URI = r.URL.Path
			trace.HttpStatus = http.StatusOK
			trace.HttpMethod = r.Method
			trace.Ingress = &storage.TrafficIngress{&requestBytes}
			db.StoreTrace(trace)
			log.Printf("[DUMP] traffic for request %s dumped with id: %s\n", trace.URI, trace.ID)

			wrapped := NewCaptureWriter(w, trace)

			// TODO error handling
			proxy.ServeHTTP(wrapped, r)
			return
		}

		log.Printf("[DUMP] could not dump traffic for request %s: %s\n", r.URL.Path, err)
		// TODO error handling
		proxy.ServeHTTP(w, r)
	})
}


// TODO : to remove if happens to be useless
const (
		NOT_STARTED		= 0
		HEADERS			= 1
		BODY			= 2
		COMPLETED		= 3
		ABORTED			= 4 // closed by the client
		TIMED_OUT       = 5 // closed before completion
)

type captureWriter struct {
	http.ResponseWriter
	trace		*storage.TrafficTrace
	state       int // NOT_STARTED => HEADERS (writing headers) => BODY (writing body) => COMPLETED
}


// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewCaptureWriter(w http.ResponseWriter, trace *storage.TrafficTrace) http.ResponseWriter {
	return captureWriter{w, trace, NOT_STARTED}
}


func (cw captureWriter) WriteHeader(status int) {
	cw.state = HEADERS
	cw.trace.HttpStatus = status
	cw.ResponseWriter.WriteHeader(status)
}


func (cw captureWriter) Write(bytes []byte) (int, error) {
	// TODO append bytes
	cw.trace.Egress = &storage.TrafficEgress{&bytes}

	// Write bytes to response
	size, err := cw.ResponseWriter.Write(bytes)
	if err != nil {
		log.Printf("[INFO] Could not write response bytes for request %s: %s\n",  cw.trace.URI, err)
		//TODO throw error
	}

	cw.trace.Length += size

	log.Printf("[DUMP] egress for %s:\n%s", cw.trace.URI, string(bytes))
	db.StoreTrace(cw.trace)


	return size, err
}


//func (cw captureWriter) CloseNotify() <-chan bool {
//	log.Printf("[DUMP] close notify for request %s\n", cw.trace.URI)
//  return cw.ResponseWriter.(http.CloseNotifier).CloseNotify()
//}



