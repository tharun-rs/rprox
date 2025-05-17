# rprox
**A simple Redis-backed reverse proxy server**  

<img src="assets/rprox_logo.png" alt="rprox logo" width="100" height="100"/>

## Features  
- Dynamic route management via **Redis**.  
- **TTL (Time-To-Live)** support for temporary routes.  
- Lightweight, Docker-ready, and easy to deploy.  



## Requirements  
- **Golang** (v1.20 or later)  
- **Docker** (with Docker Compose)  
- **Redis server** (for route storage)  



## Installation  

```sh
git clone https://github.com/tharun-rs/rprox.git
cd rprox
make all
docker-compose up -d --build
```



## Routing

`rprox` uses the **lowest-level subdomain** as the routing key.

For example, a request to:

```
<path-key>.domain.com/path
````

is routed internally to the backend registered under `/path` using the key `<path-key>`.

### How It Works

- The reverse proxy must be running on `domain.com`.
- It extracts the lowest subdomain (`<path-key>`) from the `Host` header.
- Then it matches that key with a path-to-target mapping stored in Redis.

### Example

Assume the following mapping is registered:

```bash
rproxctl -cmd set -path dev -target http://localhost:5000 -ttl 3600
````

And the proxy is running at `domain.com`. Then:

```
Request:  dev.domain.com/resource
Routes to: http://localhost:5000/resource
```

## rproxctl Usage

`rproxctl` is the command-line tool used to manage reverse proxy routes dynamically via Redis.

### General Syntax

```bash
rproxctl -cmd <command> -path <url-subdomain> [other-options]
```

Available Commands
### 1. set
Registers a new proxy path and backend target with a time-to-live (TTL).
```bash
rproxctl -cmd set -path example -target http://localhost:9000 -ttl 3600
```
- path: Path to proxy (e.g., example)
- target: Target backend URL (e.g., http://localhost:9000)
- ttl: Time in seconds before the mapping expires

### 2. get
Retrieves the current backend target mapped to a given path.
```bash
rproxctl -cmd get -path example
```
Returns: Target backend URL or an error if the path is not found

### 3. extend
Extends the TTL of an existing mapping without changing the target.
```bash
rproxctl -cmd extend -path example -ttl 1800
```
- ttl: Additional time (in seconds) to extend the TTL