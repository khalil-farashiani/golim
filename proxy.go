package main

import (
	"fmt"
	"github.com/khalil-farashiani/golim/role"
	"html"
	"io"
	"log"
	"net"
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

var proxyClient = client

func runProxy(g *golim) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := html.EscapeString(r.URL.Path)
		g.limiterRole = &limiterRole{
			operation: method,
			endPoint:  path,
		}
		currentUserRole, needToCheckRequest, err := g.getRole(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if needToCheckRequest && !isOkRequest(r, g, currentUserRole) {
			http.Error(w, slowDownError, http.StatusTooManyRequests)
			return
		}
		proxyRequest(w, r, g, currentUserRole)
	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request, g *golim, role role.GetRoleRow) {
	newURL := buildURL(role)
	proxyReq := createProxyRequest(r, newURL)
	copyHeaders(r, proxyReq)
	resp := sendProxyRequest(proxyReq)
	defer resp.Body.Close()
	copyResponseHeaders(resp, w)
	writeResponse(w, resp)
}

func buildURL(role role.GetRoleRow) url.URL {
	return url.URL{
		Scheme: "http",
		Host:   role.Destination.String,
		Path:   role.Endpoint,
	}
}

func createProxyRequest(r *http.Request, newURL url.URL) *http.Request {
	proxyReq, _ := http.NewRequest(r.Method, newURL.String(), r.Body)
	return proxyReq
}

func copyHeaders(r *http.Request, proxyReq *http.Request) {
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}
}

func sendProxyRequest(proxyReq *http.Request) *http.Response {
	resp, _ := proxyClient.Do(proxyReq)
	return resp
}

func copyResponseHeaders(resp *http.Response, w http.ResponseWriter) {
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
}

func writeResponse(w http.ResponseWriter, resp *http.Response) {
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
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
		Handler: runProxy(g),
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
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
		ip, _, err := net.SplitHostPort(ipAddress)
		if err == nil {
			ipAddress = ip
		}
	}
	return ipAddress
}
