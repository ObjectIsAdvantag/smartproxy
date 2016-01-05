package main

import (
	"log"
	"flag"
	"net/http"
	"net/http/httputil"
	"fmt"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "8080", "ip port on localhost")
	flag.Parse()

	log.Printf("[INFO] Starting Echo Service on port %s", port)

	addRoutes()

	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal("[FATAL] ",err)
	}
}


func addRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		requestBytes, err := httputil.DumpRequest(req, true)
		if err == nil {
			log.Printf(">>>> %s\n", string(requestBytes))
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		fmt.Fprint(w, `{ "state":"ok"}`)
	})
}