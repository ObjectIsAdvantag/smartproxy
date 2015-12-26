/*

Wrapper around http.ResponseWriter to capture egress traffic

Inspired from Negroni : https://github.com/codegangsta/negroni/blob/master/response_writer.go
 */
package main

import (
	"log"

	"net/http"
	"net/http/httputil"
)

func CreateRequestsDumper(proxy *httputil.ReverseProxy) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		requestBytes, err := httputil.DumpRequest(r, true)

		if err != nil {
			log.Printf("[DEBUG] Could not dump request %s: %s\n", path, err)
		} else {
			// TODO dump to memory or into some data lake
			log.Printf("[DUMP] ingress for %s\n", string(requestBytes))
		}

		wrapped := NewCaptureWriter(w, path)
		proxy.ServeHTTP(wrapped, r)
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
	path 		string
	httpStatus  int // HTTP Status code if set, 0 otherwise
	size        int // bytes written so far, 0 if not written
	state       int // NOT_STARTED => HEADERS (writing headers) => BODY (writing body) => COMPLETED
}

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface {
	http.ResponseWriter
	// Status returns the status code of the response or 0 if the response has not been written.
	HttpStatus() int
	// Size returns the size of the response body.
	Size() int
}


// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewCaptureWriter(w http.ResponseWriter, path string) ResponseWriter {
	return captureWriter{w, path, http.StatusOK, 0, NOT_STARTED}
}


func (cw captureWriter) WriteHeader(status int) {
	cw.httpStatus = status
	cw.state = HEADERS
	cw.ResponseWriter.WriteHeader(status)
}


func (cw captureWriter) Write(b []byte) (int, error) {
	// TODO dump to memory or into some data lake
	log.Printf("[DUMP] egress for %s:\n%s", cw.path, string(b))

	// Write bytes to response
	size, err := cw.ResponseWriter.Write(b)
	if err != nil {
		log.Printf("[DEBUG] Could not dump response %s: %s\n", cw.path, err)
	}

	cw.size += size
	return size, err
}


func (cw captureWriter) HttpStatus() int {
	return cw.httpStatus
}


func (cw captureWriter) Size() int {
	return cw.size
}

func (cw captureWriter) CloseNotify() <-chan bool {
	return cw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}



