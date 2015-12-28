/*

Wrapper around http.ResponseWriter to capture egress traffic

Inspired from Negroni : https://github.com/codegangsta/negroni/blob/master/response_writer.go
 */
package main

import (
	"log"
	"time"
	"strings"

	"net/http"
	"net/http/httputil"

	storage "github.com/ObjectIsAdvantag/smartproxy/storage"
)


var DB *storage.TrafficStorage = storage.OnDiskTrafficStorage()


func CreateTrafficDumper(proxy *httputil.ReverseProxy, pattern *string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestBytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("[WARNING] PROXY could not dump traffic for request %s: %s\n", r.URL.Path, err)

			// TODO error handling
			proxy.ServeHTTP(w, r)
		}

		trace := DB.CreateTrace()
		trace.URI = "/" + strings.TrimPrefix(r.URL.Path, *pattern)
		trace.HttpStatus = http.StatusOK
		trace.HttpMethod = r.Method
		trace.Ingress = &storage.TrafficIngress{&requestBytes}
		log.Printf("[DEBUG] PROXY traffic for request %s dumped with id: %s\n", trace.URI, trace.ID)

		wrapped := NewCaptureWriter(w, trace)

		trace.Start = time.Now()

		// TODO error handling
		proxy.ServeHTTP(wrapped, r)

		trace.End = time.Now()
		DB.StoreTrace(trace)

		return
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

	// TODO append bytes if called several time
	cw.trace.Egress = &storage.TrafficEgress{&bytes}

	// Write bytes to response
	size, err := cw.ResponseWriter.Write(bytes)
	if err != nil {
		log.Printf("[WARNING] PROXY Could not write response bytes for request %s: %s\n",  cw.trace.URI, err)
		//TODO throw error
	}

	cw.trace.Length += size
	cw.trace.End = time.Now()

	//log.Printf("[DEBUG] PROXY egress for %s:\n%s", cw.trace.URI, string(bytes))
	//DB.StoreTrace(cw.trace)

	return size, err
}


//func (cw captureWriter) CloseNotify() <-chan bool {
//	log.Printf("[DUMP] close notify for request %s\n", cw.trace.URI)
//  return cw.ResponseWriter.(http.CloseNotifier).CloseNotify()
//}



