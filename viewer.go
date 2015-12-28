/*

Viewer to inspect traffic

 */
package main

import (
	"log"
	"fmt"

	"net/http"
)

func AddTrafficViewer(route string) {

	http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[DEBUG] VIEWER ---> %s %s", req.Method, req.URL.String())

		if !authorizeOnlyGET(w, req) {
			return
		}

		if isHuman() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, "<html><head><title>Traffic inspection</title></head><body><h1>Traffic inspection</h1>");

			DB.DisplayLatestTraces(w, 10)

			fmt.Fprint(w, "</body></html>")
			return
		}

		// Not implemented
		log.Printf("[DEBUG] VIEWER Traffic inspection for machine is not implemented yet, try same URI from a WebBrowser")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Not implemented")
	})

	http.HandleFunc(route+"/last", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[DEBUG] VIEWER ---> %s %s", req.Method, req.URL.String())

		if !authorizeOnlyGET(w, req) {
			return
		}

		if isHuman() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, "<html><head><title>Traffic inspector</title></head><body><h1>Last capture</h1>");
			DB.DisplayLastTrace(w)
			fmt.Fprint(w, "</body></html>")

			return
		}

		// Not implemented
		log.Printf("[DEBUG] VIEWER Traffic inspection for machine is not implemented yet, try same URI from a WebBrowser")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Not implemented")
	})

	http.HandleFunc(route+"/first", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[DEBUG] VIEWER ---> %s %s", req.Method, req.URL.String())

		if !authorizeOnlyGET(w, req) {
			return
		}

		if isHuman() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, "<html><head><title>Traffic inspector</title></head><body><h1>Last capture</h1>");
			DB.DisplayFirstTrace(w)
			fmt.Fprint(w, "</body></html>")

			return
		}

		// Not implemented
		log.Printf("[DEBUG] VIEWER Traffic inspection for machine is not implemented yet, try same URI from a WebBrowser")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Not implemented")
	})

	http.HandleFunc(route+"/next", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[DEBUG] VIEWER ---> %s %s", req.Method, req.URL.String())

		if !authorizeOnlyGET(w, req) {
			return
		}

		if isHuman() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, "<html><head><title>Traffic inspector</title></head><body><h1>Last capture</h1>");
			DB.DisplayNextTrace(w)
			fmt.Fprint(w, "</body></html>")

			return
		}

		// Not implemented
		log.Printf("[DEBUG] VIEWER Traffic inspection for machine is not implemented yet, try same URI from a WebBrowser")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Not implemented")
	})

	http.HandleFunc(route+"/prev", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[DEBUG] VIEWER ---> %s %s", req.Method, req.URL.String())

		if !authorizeOnlyGET(w, req) {
			return
		}

		if isHuman() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, "<html><head><title>Traffic inspector</title></head><body><h1>Last capture</h1>");
			DB.DisplayPrevTrace(w)
			fmt.Fprint(w, "</body></html>")

			return
		}

		// Not implemented
		log.Printf("[DEBUG] VIEWER Traffic inspection for machine is not implemented yet, try same URI from a WebBrowser")
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Not implemented")
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
