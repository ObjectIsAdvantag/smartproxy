/*
Adaptation of http/httputil/NewSingleHostReverseProxy with no URL rewrite,
ie, the proxy route/pattern is not prefixed to incoming requests URLs
 */
package main

import (
	"net/url"
	"net/http"
	"net/http/httputil"
	"strings"
)


func TransparentReverseProxy(target *url.URL, pattern *string) *httputil.ReverseProxy {

		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			if *pattern != "/" {
				req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, *pattern)
			}
		}
		return &httputil.ReverseProxy{Director: director}
}
