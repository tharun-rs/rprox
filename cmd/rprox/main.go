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
	host := r.Host

	targetURL, err := store.Get(host)
	if err != nil || targetURL == "" {
		http.Error(w, "No route found", http.StatusBadGateway)
		return
	}

	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
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
