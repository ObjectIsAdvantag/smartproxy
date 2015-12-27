/*

SmartProxy is a companion for Web API developers by mocking, filtering, replaying, "diff'ing" HTTP req/responses
It may also help Web API hosting via diff'ing between API versions and taking actions when errors.

SmartProxy acts as a reverse proxy that

- maintains an history of request
- allows to inspect them
- allows to modify them : YOU take action

 *  Inspired by eBay/fabio, mitmproxy, WireMock, mashape/kong
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

    // _ "expvar" // adds a memstats endpoint handler at /debug/vars
	//_ "net/http/pprof"
)


// version contains the version number
// It is set by build/release.sh for tagged releases
// so that 'go get' just works.
//
// It is also set by the linker when fabio
// is built via the Makefile or the build/docker.sh
// script to ensure the correct version nubmer
var version = "v0.3"


func main() {
	var showVersion, dumpRequest bool
	var name, port, serve, route, healthcheck string
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&dumpRequest, "dump", false, "dumps ingoing and outgoing proxied traffic")
	flag.StringVar(&serve, "serve", "127.0.0.1:8080", "host or host:port of the proxied service, defaults to 127.0.0.1:8080")
	// WORKAROUND do not use absolute path (starting with /), because go runtime expands as a directory
	flag.StringVar(&route, "route", "", "relative path to the proxied service, defaults to /, on WINDOWS : do not prefix with /")
	flag.StringVar(&port, "port", "9090", "ip port of reverse proxy, defaults to 9090")
	flag.StringVar(&name, "name", "SmartProxy", "name of the service, defaults to SmartProxy")
	flag.StringVar(&healthcheck, "healthcheck", "/ping", "healthcheck endpoint, defaults to /ping, on WINDOWS : do not prefix with /")
	flag.Parse()

	if showVersion {
		fmt.Printf("SmartProxy version %s", version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("[FATAL] Invalid port: %s (%s)\n", port, err)
	}

	log.Printf("[INFO] Starting %s, version: %s", name, version)

	// start http server
	go func() {
		// register reverse proxy and middlewares
		endpoint := &url.URL{Scheme:"http", Host:serve}
		pattern := computeProxyPath(route)
		proxy := CreateReverseProxy(endpoint, &pattern) // *ReverseProxy
		if dumpRequest {
			handler := CreateCaptureMiddleware(proxy)
			http.Handle(pattern, handler)
		} else {
			http.HandleFunc(pattern, proxy.ServeHTTP)
		}

		// register health check
		ping := computeHealthcheckPath(healthcheck)
		http.HandleFunc(ping, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[INFO] Health check")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(w, `{ "version":"%s", "state":"active", "name":"%s", "port":"%s", "serving":"http://%s", "via":"%s", "healthcheck":"%s", "dump":"%v"}`, version, name, port, serve, pattern, healthcheck, dumpRequest)
		})

		// add a default route if the proxy is not registered on /
		if pattern != "/" {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				log.Printf("[INFO] No route registered for %s", r.RequestURI)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{ "state":"error", "description":"nothing here" }`)
			})
		}

		log.Printf("[INFO] Listening on http://localhost:%s\n", port)
		log.Printf("[INFO] Serving http://%s via path %s\n", serve, pattern)
		if err := http.ListenAndServe(":" + port, nil); err != nil {
			log.Fatal("[FATAL] ",err)
		}
	}()

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}

func computeProxyPath(route string) string {
	// Ensure a leading and ending / to the route
	//    - leading / is due to an issue on windows with Args, which requires to pass the route arg as relative (no leading /)
	//           see https://golang.org/src/os/proc.go, line 19
	//           see https://golang.org/src/syscall/exec_windows.go, line 156
	//    - ending / is required to serve all traffic to the proxy /route, /route/ and /route/*, and not only the /route URL
	pattern := route
	if !(strings.HasPrefix(pattern, "/")) {
		pattern = "/" + pattern
	}
	if !(strings.HasSuffix(pattern, "/")) {
		pattern = pattern + "/"
	}
	return pattern
}

func computeHealthcheckPath(route string) string {
	// Ensure a leading / to the route
	pattern := route
	if !(strings.HasPrefix(pattern, "/")) {
		pattern = "/" + pattern
	}
	return pattern
}

