/*
Adaptation of http/httputil/NewSingleHostReverseProxy with no URL rewrite,
ie, the proxy route/pattern is not prefixed to incoming requests URLs
 */
package main

import (
	"strings"

	"net/url"
	"net/http"
	"net/http/httputil"
	"log"
)


func CreateReverseProxy(target *url.URL, pattern *string) *httputil.ReverseProxy {

		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			if *pattern != "/" {
				req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, *pattern)
			}
		}
		return &httputil.ReverseProxy{Director: director}
}

func RegisterMiddleware(proxy *httputil.ReverseProxy, dump bool) http.Handler {
    if !dump {
		return proxy
	}

	log.Printf("[INFO] Registering middleware to dump traffic")
	return addRequestsDumper(proxy)
}

func addRequestsDumper(proxy *httputil.ReverseProxy) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = "breakpoint"

		path := r.URL.Path
		requestBytes, err := httputil.DumpRequest(r, true)

		if err != nil {
			log.Printf("[DEBUG] Could not dump request %s: %s\n", path, err)
		} else {
			// TODO dump to memory or into some data lake
			log.Printf("[DUMP] ingress for %s\n", string(requestBytes))
		}

		wrapped := NewCaptureWriter(w, path)
		proxy.ServeHTTP(wrapped, r)
	})
}

