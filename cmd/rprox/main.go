package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/tharun-rs/rprox/config"
	"github.com/tharun-rs/rprox/logger"
	"github.com/tharun-rs/rprox/redis"
)

var store *redis.RedisClient

func main() {
	config.Init()
	store = &redis.RedisClient{}

	if err := store.Init(config.Cfg); err != nil {
		logger.Log.Errorf("Failed to connect to Redis: %v", err)
		return
	}

	logger.Log.Infof("Redis initialized: %s", config.Cfg.RedisURL)

	http.HandleFunc("/", handleReverseProxy)

	addr := config.Cfg.AppPort
	logger.Log.Infof("Reverse proxy listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Log.Errorf("Failed to start server: %v", err)
	}
}

func handleReverseProxy(w http.ResponseWriter, r *http.Request) {
	// Extract the first path segment as the key for Redis
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathSegments) == 0 || pathSegments[0] == "" {
		http.Error(w, "No route found", http.StatusBadGateway)
		return
	}

	key := pathSegments[0]

	targetURLStr, err := store.Get(key)
	if err != nil || targetURLStr == "" {
		http.Error(w, "No route found", http.StatusBadGateway)
		return
	}

	target, err := url.Parse(targetURLStr)
	if err != nil {
		http.Error(w, "Invalid target", http.StatusInternalServerError)
		return
	}

	// Build the new path by appending the remaining path segments after the first
	remainingPath := "/"
	if len(pathSegments) > 1 {
		remainingPath += strings.Join(pathSegments[1:], "/")
	}

	// Preserve query parameters as well
	if r.URL.RawQuery != "" {
		remainingPath += "?" + r.URL.RawQuery
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Set Host header to target host
		req.Host = target.Host

		// Override the request path with the target's base path + remaining path from original URL
		req.URL.Path = singleJoiningSlash(target.Path, remainingPath)
		req.URL.RawPath = ""

		// Optionally update req.URL.RawQuery if needed (already preserved above)
	}

	proxy.ServeHTTP(w, r)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	default:
		return a + b
	}
}
