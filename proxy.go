package main

import (
	"fmt"
	"github.com/khalil-farashiani/golim/role"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var customTransport = &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
}

var client = &http.Client{
	Timeout:   time.Second * 10,
	Transport: customTransport,
}

func runProxy(g *golim) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := html.EscapeString(r.URL.Path)
		g.limiterRole = &limiterRole{
			operation: method,
			endPoint:  path,
		}
		role, needToCheckRequest, err := g.getRole(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if needToCheckRequest && !isOkRequest(r, g, role) {
			http.Error(w, slowDownError, http.StatusTooManyRequests)
			return
		}
		newURL := url.URL{
			Scheme: "http",
			Host:   role.Destination.String,
			Path:   role.Endpoint,
		}
		targetURL := newURL
		proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			http.Error(w, createProxyError, http.StatusInternalServerError)
			return
		}
		for name, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, sendingProxyError, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func isOkRequest(r *http.Request, g *golim, role role.GetRoleRow) bool {
	ctx := r.Context()
	userIP := readUserIP(r)
	capacity := g.cache.getUserRequestCap(ctx, userIP, g, role)
	if capacity > 0 {
		go g.cache.decreaseCap(r.Context(), userIP, g.limiterRole)
		return true
	}
	return false
}

func startServer(g *golim) (interface{}, error) {
	portStr := fmt.Sprintf(":%d", g.port)
	server := http.Server{
		Addr:    portStr,
		Handler: http.HandlerFunc(runProxy(g)),
	}

	// Start the server and log any errors
	log.Printf("Starting golim on %d", g.port)
	err := server.ListenAndServe()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
