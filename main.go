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
	// WORKAROUND do not use absolute path (starting with /), because go runtime expands as a directory
	flag.StringVar(&route, "route", "", "relative path to the proxied service, defaults to /, on WINDOWS : no not prefix with /")
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
		_="breakpoint"

		// register reverse proxy
		endpoint := &url.URL{Scheme:"http", Host:serve}
		proxy := httputil.NewSingleHostReverseProxy(endpoint) // *ReverseProxy
		pattern := computeProxyPath(route)
		http.HandleFunc(pattern, proxy.ServeHTTP)

		// register health check, except if proxy is not registered on /health
		// TODO : check if the / suffix does not do the trick
		if pattern != "/health/" {
			http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				fmt.Fprintf(w, `{ "state":"active", "name":"%s", "port":"%s", "serving":"http://%s%s"}`, name, port, serve, pattern)
			})
		}

		// add a default route if the proxy is not registered on /
		if pattern != "/" {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				log.Printf("[INFO] No route registered for %s", r.RequestURI)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{ "state":"error", "description":"nothing here" }`)
			})
		}

		log.Printf("[INFO] Listening on port %s, serving http://%s%s", port, serve, pattern)
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
	// Amend a starting and ending / to the route
	//    - starting / is due to an issue on windows with Args, which requires to pass the route arg as relative (no leading /)
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
