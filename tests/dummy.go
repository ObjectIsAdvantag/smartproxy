package main

import (
	"log"
	"flag"
	"net/http"
	"fmt"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "8080", "ip port on localhost")
	flag.Parse()

	log.Printf("[INFO] Starting dummy server on port %s", port)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, "{ \"state\":\"ok\" }")
	})

	http.ListenAndServe(":"+port, nil)
}
