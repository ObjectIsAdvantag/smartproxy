package main

import (
	"log"
	"flag"
	"net/http"
	"fmt"
	"encoding/json"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "8080", "ip port on localhost")
	flag.Parse()

	log.Printf("[INFO] Starting dummy server on port %s", port)

	addRoutes()
	addRoutesWithoutContentTypes() // automatically computed by golang

	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal("[FATAL] ",err)
	}
}


func addRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprint(w, "<html><body>Hello Stève</body></html>")
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		fmt.Fprint(w, `{ "state":"ok"}`)
	})

	http.HandleFunc("/json/encoded", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		type Message struct {
			State   string `json:"state"`
			Encoded bool `json:"encoded"`
		}
		enc := json.NewEncoder(w)
		m := &Message{State:"ok", Encoded:true}
		enc.Encode(m)
	})

	http.HandleFunc("/txt", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		fmt.Fprint(w, "This is an UTF-8 message ! Stève ?")
	})
}


func addRoutesWithoutContentTypes() {

	http.HandleFunc("/untyped", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		fmt.Fprint(w, "<html><body>Hello Stève</body></html>")
	})

	http.HandleFunc("/untyped/json", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		fmt.Fprint(w, `{ "state":"ok"}`)
	})

	http.HandleFunc("/untyped/json/encoded", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		type Message struct {
			State string `json:"state"`
			Encoded bool `json:"encoded"`
		}
		enc := json.NewEncoder(w)
		m := &Message{State:"ok", Encoded:true}
		enc.Encode(m)
	})

	http.HandleFunc("/untyped/txt", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("---> %s %s", req.Method, req.URL.String())

		fmt.Fprint(w, "This is an UTF-8 message ! Stève ?")
	})
}