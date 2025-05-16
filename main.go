package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"

    // "github.com/tharun-rs/rprox/logger"
)

func main() {
    targetURL, err := url.Parse("http://httpbin.org")
    if err != nil {
        log.Fatalf("Error parsing target URL: %v", err)
    }

    proxy := httputil.NewSingleHostReverseProxy(targetURL)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        proxy.ServeHTTP(w, r)
    })

    log.Println("Starting reverse proxy on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
