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
	log.Printf("[INFO] Creating Reverse Proxy at %s\n", *pattern)
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		if *pattern != "/" {
			req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, *pattern)
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func CreateCaptureMiddleware(proxy *httputil.ReverseProxy) http.Handler {
	log.Printf("[INFO] Adding Middleware to capture traffic\n")
	return CreateTrafficDumper(proxy)
}




