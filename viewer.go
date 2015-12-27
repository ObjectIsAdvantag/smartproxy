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
		_ = "breakpoint"
		log.Printf("[VIEWER] ---> %s %s", req.Method, req.URL.String())

		if (req.Method != "GET") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, "only get is accepted here</body></html>")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<html><head><title>Traffic inspection</title></head><body><h1>Traffic inspection</h1>");

		DB.DisplayLastTraces(w, 10)

		fmt.Fprint(w, "</body></html>")
	})

	http.HandleFunc(route+"/", func(w http.ResponseWriter, req *http.Request) {
		_ = "breakpoint"
		log.Printf("[VIEWER] ---> %s %s", req.Method, req.URL.String())

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		fmt.Fprint(w, `{ "state":"ok"}`)
	})


}


