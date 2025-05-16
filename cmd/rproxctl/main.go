package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tharun-rs/rprox/config"
	"github.com/tharun-rs/rprox/logger"
	"github.com/tharun-rs/rprox/redis"
)

func main() {
	config.Init()

	cmd := flag.String("cmd", "", "Command to run: set|get|extend")
	path := flag.String("path", "", "Route path (e.g., /api)")
	target := flag.String("target", "", "Target URL (for set)")
	ttl := flag.String("ttl", "0", "TTL in seconds (for set or extend), 0 means no expiration")
	flag.Parse()

	if *cmd == "" || *path == "" {
		fmt.Println("Usage:")
		fmt.Println("  rproxctl -cmd set -path /example -target http://localhost:9000 -ttl 3600")
		fmt.Println("  rproxctl -cmd get -path /example")
		fmt.Println("  rproxctl -cmd extend -path /example -ttl 1800")
		os.Exit(1)
	}

	client := redis.RedisClient{}
	if err := client.Init(config.Cfg); err != nil {
		logger.Log.Errorf("Failed to connect to Redis: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	switch *cmd {
	case "set":
		if *target == "" {
			logger.Log.Errorf("Target is required for set")
			os.Exit(1)
		}

		ttlSeconds, err := strconv.Atoi(*ttl)
		if err != nil || ttlSeconds < 0 {
			logger.Log.Errorf("Invalid TTL value: %s", *ttl)
			os.Exit(1)
		}
		var expiration time.Duration
		if ttlSeconds > 0 {
			expiration = time.Duration(ttlSeconds) * time.Second
		} else {
			expiration = 0
		}

		err = client.Put(*path, *target, expiration)
		if err != nil {
			logger.Log.Errorf("Failed to set route: %v", err)
			os.Exit(1)
		}
		logger.Log.Infof("Set route: %s => %s with TTL %v", *path, *target, expiration)

	case "get":
		val, err := client.Get(*path)
		if err != nil {
			logger.Log.Errorf("Failed to get route: %v", err)
			os.Exit(1)
		}
		logger.Log.Infof("Route: %s => %s", *path, val)

	case "extend":
		ttlSeconds, err := strconv.Atoi(*ttl)
		if err != nil || ttlSeconds <= 0 {
			logger.Log.Errorf("Invalid TTL value for extend: %s", *ttl)
			os.Exit(1)
		}
		expiration := time.Duration(ttlSeconds) * time.Second

		err = client.Extend(*path, expiration)
		if err != nil {
			logger.Log.Errorf("Failed to extend TTL: %v", err)
			os.Exit(1)
		}
		logger.Log.Infof("Extended TTL for %s by %v", *path, expiration)

	default:
		logger.Log.Errorf("Unknown command: %s", *cmd)
		os.Exit(1)
	}
}
