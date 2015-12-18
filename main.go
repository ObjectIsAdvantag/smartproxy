/*

SmartProxy is a companion for Web API developers by mocking, filtering, replaying, "diff'ing" HTTP req/responses
It may also help Web API hosting via diff'ing between API versions and taking actions when errors.

SmartProxy acts as a reverse proxy that

- maintains an history of request
- allows to inspect them
- allows to modify them : YOU take action

 *  Inspired by eBay/fabio, goproxy
 */
package main

import (
	"log"
	"fmt"

	"flag"
	"os"
	"os/signal"
	"strconv"

	"net/url"
	"net/http"
	"net/http/httputil"
)


// version contains the version number
// It is set by build/release.sh for tagged releases
// so that 'go get' just works.
//
// It is also set by the linker when fabio
// is built via the Makefile or the build/docker.sh
// script to ensure the correct version nubmer
var version = "draft"




func main() {

	var v bool
	var name, port, serve string
	flag.BoolVar(&v, "v", false, "show version")
	flag.StringVar(&serve, "ip", "127.0.0.1:8080", "address of the proxied service")
	flag.StringVar(&port, "port", "9090", "ip port of reverse proxy")
	flag.StringVar(&name, "name", "SmartProxy", "name of the service")
	flag.Parse()

	if v {
		fmt.Printf("SmartProxy version %s, build undefined", version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("[FATAL] Invalid port: %s (%s)\n", port, err)
	}

	log.Printf("[INFO] Starting version %s of %s", version, name)

	// start http server
	go func() {
		log.Printf("[INFO] Listening on port %s, serving %s", port, serve)
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   serve,
		})
		if err := http.ListenAndServe(":" + port, proxy); err != nil {
			log.Fatal("[FATAL] ",err)
		}
	}()

	// register health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		//w.Write(byte[]("{ 'state':'ok', 'name':'smart proxy', 'port':'9090', 'serving':'127.0.0.1:8080'  }"))
		fmt.Fprintf(w, "{ 'state':'ok', 'name':'smart proxy', 'port':'9090', 'serving':'127.0.0.1:8080'  }")
	})

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}

