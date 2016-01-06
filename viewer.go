/*

Viewer to inspect traffic

 */
package main

import (
	"log"
	"fmt"
	"strings"

	"net/http"
)

func AddTrafficViewer(route string) {

	http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[DEBUG] VIEWER ---> %s %s", req.Method, req.URL.String())

		if !authorizeOnlyGET(w, req) {
			return
		}

		if isHuman() {
			DB.DisplayLatestTraces(w, route, 20)
			return
		}

		// Not implemented
		log.Printf("[DEBUG] VIEWER Traffic inspection for machine is not implemented yet, try same URI from a WebBrowser")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Not implemented")
	})

	http.HandleFunc(route+"/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[DEBUG] VIEWER ---> %s %s", req.Method, req.URL.String())

		if !authorizeOnlyGET(w, req) {
			return
		}

		id := strings.TrimPrefix(req.URL.Path, route+"/")
		if id == "" {
			DB.DisplayLatestTraces(w, route, 20)
			return
		}

		log.Printf("[DEBUG] VIEWER inspect trace with id %s\n", id)
		DB.DisplayTraceDetails(w, route, id)
	})

}

func isHuman() bool {
	return true
}

// Returns true if the req method is authorized
func authorizeOnlyGET (w http.ResponseWriter, req *http.Request) bool {

	if req.Method == "GET" {
		return true
	}

	if isHuman() {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<hmtl><body><h1>Error</h1><p>only GET is accepted here</p></body></html>")
		return false
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{ "error":"404", "description":"only GET method is accepted here"}`)
	return false

}
