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
	"strings"

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
	var name, port, serve, route string
	flag.BoolVar(&v, "v", false, "show version")
	flag.StringVar(&serve, "serve", "127.0.0.1:8080", "host or host:port of the proxied service")
	flag.StringVar(&route, "route", "/proxy", "reverse proxy path to join the proxied service")
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
		// register reverse proxy
		endpoint := &url.URL{
			Scheme: "http",
			Host:   serve,
		}
		proxy := httputil.NewSingleHostReverseProxy(endpoint) // *ReverseProxy
		// Make sure the pattern ends with a / to serve all traffic to the proxy
		if !strings.HasSuffix(route, "/") {
			route = route + "/"
		}
		http.HandleFunc(route, proxy.ServeHTTP)

		// register health check endpoint
		//TODO : reject route == "/health"
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "{ 'state':'active', 'name':%s, 'port':%s, 'serving':%s  }", name, port, serve)
		})

		// add a default route if reverse proxy does not register on /
		if route != "/" {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				log.Printf("[INFO] No route registered for %s", r.RequestURI)
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "{ 'state':'error', 'description':'nothing here' }")
			})
		}

		log.Printf("[INFO] Listening on port %s, serving %s", port, serve)
		if err := http.ListenAndServe(":" + port, nil); err != nil {
			log.Fatal("[FATAL] ",err)
		}
	}()

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}

