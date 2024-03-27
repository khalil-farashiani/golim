package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
)

var customTransport = http.DefaultTransport

func runProxy(g *golim) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := html.EscapeString(r.URL.Path)

		g.limiterRole = &limiterRole{
			operation: method,
			endPoint:  path,
		}
		if !isOkRequest(r, g) {
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}

		role, err := g.getRole(r.Context())
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

func isOkRequest(r *http.Request, g *golim) bool {
	ctx := r.Context()
	userIP := readUserIP(r)
	capacity := g.cache.getUserRequestCap(ctx, userIP, g.limiterRole)
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
