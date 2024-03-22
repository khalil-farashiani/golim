package main

import (
	"database/sql"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
)

var operationStringToID = map[string]int{
	"GET":    1,
	"POST":   2,
	"PUT":    3,
	"PATCH":  4,
	"DELETE": 5,
}

var customTransport = http.DefaultTransport

func runProxy(g *golim, db *sql.DB, cache *cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := html.EscapeString(r.URL.Path)
		lr := &limiterRole{
			operation: operationStringToID[method],
			endPoint:  path,
		}
		g.limiterRole = lr
		role, err := g.getRole(r.Context(), db, cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isOkRequest(r, lr, cache) {
			http.Error(w, fmt.Sprintf("slow down"), http.StatusTooManyRequests)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
			return
		}
		for name, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		resp, err := customTransport.RoundTrip(proxyReq)
		if err != nil {
			http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
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

func isOkRequest(r *http.Request, rl *limiterRole, cache *cache) bool {
	ctx := r.Context()
	userIP := readUserIP(r)
	capacity := cache.getUserRequestCap(ctx, userIP, rl)
	if capacity > 0 {
		return true
	}
	return false
}

func startServer(g *golim, db *sql.DB, cache *cache) (interface{}, error) {
	portStr := fmt.Sprintf(":%d", g.port)
	server := http.Server{
		Addr:    portStr,
		Handler: http.HandlerFunc(runProxy(g, db, cache)),
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
