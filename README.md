# StremThru

<div align="center">

**Privacy-First Debrid Proxy for Stremio**

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://hub.docker.com/r/muniftanjim/stremthru)

*A secure proxy service that shields your IP address from debrid providers like RealDebrid, AllDebrid, and more.*

</div>

---

## Table of Contents

- [Overview](#overview)
  - [Key Features](#key-features)
  - [Supported Debrid Services](#supported-debrid-services)
  - [Architecture Summary](#architecture-summary)
- [Quick Start](#quick-start)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Local Development](#local-development)
  - [Docker Deployment](#docker-deployment)
  - [Heroku Deployment](#heroku-deployment)
- [Configuration Reference](#configuration-reference)
  - [Environment Variables](#environment-variables)
  - [Feature Flags](#feature-flags)
  - [Tunnel Configuration](#tunnel-configuration)
- [Cloudflare WARP Integration](#cloudflare-warp-integration)
  - [Overview](#warp-overview)
  - [Setup Guide](#warp-setup-guide)
  - [Configuration](#warp-configuration)
  - [Troubleshooting WARP](#warp-troubleshooting)
- [API Reference](#api-reference)
  - [Health Endpoints](#health-endpoints)
  - [Proxy Endpoints](#proxy-endpoints)
  - [Store Endpoints](#store-endpoints)
  - [Stremio Addons](#stremio-addons)
- [System Design](#system-design)
  - [Architecture Diagram](#architecture-diagram)
  - [Request Flow](#request-flow)
  - [Data Flow](#data-flow)
  - [Component Interactions](#component-interactions)
- [Commands Reference](#commands-reference)
  - [Go Commands](#go-commands)
  - [Make Commands](#make-commands)
  - [Node/PNPM Commands](#nodepnpm-commands)
  - [Heroku Commands](#heroku-commands)
- [Operations](#operations)
  - [Logging](#logging)
  - [Monitoring](#monitoring)
  - [Database Management](#database-management)
  - [Scaling Considerations](#scaling-considerations)
- [Security](#security)
  - [Authentication](#authentication)
  - [IP Privacy Model](#ip-privacy-model)
  - [Secrets Handling](#secrets-handling)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
  - [Repository Structure](#repository-structure)
  - [Coding Standards](#coding-standards)
  - [Pull Request Process](#pull-request-process)
- [License](#license)

---

## Overview

StremThru is a Go-based proxy service designed to act as an intermediary between you and debrid services (RealDebrid, AllDebrid, Premiumize, etc.). Its primary purpose is to **protect your IP address** from being exposed to debrid providers.

### Key Features

| Feature | Description |
|---------|-------------|
| 🔒 **IP Privacy** | Debrid services only see StremThru's IP, never yours |
| 🌐 **Cloudflare WARP** | Optional integration to use Cloudflare's network |
| 📺 **Stremio Integration** | Native Stremio addon support for seamless streaming |
| 🔄 **Content Proxy** | Proxies video content with byte-range support |
| 🗄️ **Multi-Database** | Supports SQLite (default) and PostgreSQL |
| 🐳 **Docker Ready** | Multi-stage Dockerfile for minimal image size |
| ☁️ **Heroku Ready** | One-click deployment with `app.json` |

### Supported Debrid Services

| Service | Status | Code |
|---------|--------|------|
| RealDebrid | ✅ Fully Supported | `realdebrid` |
| AllDebrid | ✅ Fully Supported | `alldebrid` |
| Premiumize | ✅ Fully Supported | `premiumize` |
| DebridLink | ✅ Fully Supported | `debridlink` |
| TorBox | ✅ Fully Supported | `torbox` |
| Offcloud | ✅ Fully Supported | `offcloud` |
| EasyDebrid | ✅ Fully Supported | `easydebrid` |
| PikPak | ✅ Fully Supported | `pikpak` |

### Architecture Summary

```mermaid
graph LR
    subgraph "Your Network"
        A[Stremio Client]
    end
    
    subgraph "StremThru Server"
        B[StremThru Proxy]
        C[WireProxy/WARP]
    end
    
    subgraph "External Services"
        D[RealDebrid API]
        E[RealDebrid CDN]
    end
    
    A -->|Your IP Hidden| B
    B -->|SOCKS5| C
    C -->|Cloudflare IP| D
    C -->|Cloudflare IP| E
    
    style A fill:#e1f5fe
    style B fill:#fff3e0
    style C fill:#f3e5f5
    style D fill:#e8f5e9
    style E fill:#e8f5e9
```

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/MunifTanjim/stremthru.git
cd stremthru

# Copy environment file
cp .env.example .env

# Edit with your credentials
# STREMTHRU_PROXY_AUTH=admin:yourpassword
# STREMTHRU_STORE_AUTH=*:realdebrid:YOUR_API_TOKEN

# Run with Docker
docker-compose up -d

# Or run locally
go run --tags "fts5" .
```

---

## Installation

### Prerequisites

| Requirement | Version | Purpose |
|-------------|---------|---------|
| Go | 1.25+ | Backend compilation |
| Node.js | 22+ | Dashboard frontend |
| pnpm | 10.17+ | Package manager |
| Docker | 24+ | Container deployment (optional) |
| Heroku CLI | Latest | Heroku deployment (optional) |

### Local Development

```bash
# 1. Clone the repository
git clone https://github.com/MunifTanjim/stremthru.git
cd stremthru

# 2. Install frontend dependencies
corepack enable
pnpm install

# 3. Build the dashboard
pnpm run dash:build

# 4. Create environment file
cp .env.example .env

# 5. Configure your credentials (edit .env)
STREMTHRU_PROXY_AUTH=admin:password123
STREMTHRU_STORE_AUTH=*:realdebrid:YOUR_RD_API_KEY
STREMTHRU_DATABASE_URI=sqlite://./data/stremthru.db

# 6. Run the application
go run --tags "fts5" .

# Server starts at http://localhost:8080
```

### Docker Deployment

#### Using Docker Compose (Recommended)

```yaml
# compose.yaml
services:
  stremthru:
    build: .
    ports:
      - "8080:8080"
    environment:
      - STREMTHRU_PROXY_AUTH=admin:password123
      - STREMTHRU_STORE_AUTH=*:realdebrid:YOUR_API_KEY
      - STREMTHRU_DATABASE_URI=sqlite://./data/stremthru.db
      - STREMTHRU_STORE_CONTENT_PROXY=*:true
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f stremthru

# Stop the service
docker-compose down
```

#### Using Docker CLI

```bash
# Build the image
docker build -t stremthru .

# Run the container
docker run -d \
  --name stremthru \
  -p 8080:8080 \
  -e STREMTHRU_PROXY_AUTH=admin:password123 \
  -e STREMTHRU_STORE_AUTH=*:realdebrid:YOUR_API_KEY \
  -v $(pwd)/data:/app/data \
  stremthru
```

### Heroku Deployment

#### One-Click Deploy

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

#### Manual Deployment

```bash
# 1. Create Heroku app
heroku create your-app-name --stack container

# 2. Set environment variables
heroku config:set STREMTHRU_BASE_URL=https://your-app-name.herokuapp.com
heroku config:set STREMTHRU_PROXY_AUTH=admin:yourpassword
heroku config:set STREMTHRU_STORE_AUTH=*:realdebrid:YOUR_API_KEY
heroku config:set STREMTHRU_STORE_CONTENT_PROXY=*:true
heroku config:set STREMTHRU_STORE_TUNNEL=*:true
heroku config:set STREMTHRU_DATABASE_URI=sqlite://./data/stremthru.db

# 3. Deploy
git push heroku main

# 4. Verify deployment
heroku logs --tail
curl https://your-app-name.herokuapp.com/v0/health
```

---

## Configuration Reference

### Environment Variables

#### Core Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `STREMTHRU_ENV` | No | `dev` | Environment: `dev`, `prod`, `test` |
| `STREMTHRU_BASE_URL` | No | Auto-detected | Public URL of your instance |
| `STREMTHRU_PORT` | No | `8080` | HTTP server port |
| `STREMTHRU_DATABASE_URI` | No | `sqlite://./data/stremthru.db` | Database connection string |

#### Authentication

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `STREMTHRU_PROXY_AUTH` | **Yes** | - | Proxy auth in `user:pass` format. Multiple: `user1:pass1,user2:pass2` |
| `STREMTHRU_STORE_AUTH` | **Yes** | - | Store auth in `user:store:token` format. Use `*` for any user |
| `STREMTHRU_ADMIN_PASSWORD` | No | - | Admin password for dashboard access |

**Example:**
```bash
# Single user
STREMTHRU_PROXY_AUTH=admin:secretpassword

# Multiple users
STREMTHRU_PROXY_AUTH=user1:pass1,user2:pass2

# Store auth for any user with RealDebrid
STREMTHRU_STORE_AUTH=*:realdebrid:YOUR_RD_API_KEY

# Store auth for specific user
STREMTHRU_STORE_AUTH=john:realdebrid:JOHNS_API_KEY
```

#### Proxy & Tunnel

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `STREMTHRU_HTTP_PROXY` | No | - | HTTP/SOCKS5 proxy URL |
| `STREMTHRU_TUNNEL` | No | - | Tunnel configuration per hostname |
| `STREMTHRU_STORE_TUNNEL` | No | `*:true` | Whether to tunnel store traffic |
| `STREMTHRU_STORE_CONTENT_PROXY` | No | - | Proxy content through StremThru |

**Tunnel Configuration Syntax:**
```bash
# Tunnel all traffic
STREMTHRU_STORE_TUNNEL=*:true

# Tunnel only API calls, not streams
STREMTHRU_STORE_TUNNEL=*:api

# Disable tunnel for specific store
STREMTHRU_STORE_TUNNEL=*:true,realdebrid:false
```

#### Database

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `STREMTHRU_DATABASE_URI` | No | `sqlite://./data/stremthru.db` | Database URI |
| `DATABASE_URL` | No | - | Alternative database URI (Heroku) |

**Database URI Formats:**
```bash
# SQLite (default)
STREMTHRU_DATABASE_URI=sqlite://./data/stremthru.db

# PostgreSQL
STREMTHRU_DATABASE_URI=postgres://user:pass@host:5432/dbname
```

### Feature Flags

Control which features are enabled using `STREMTHRU_FEATURE`:

```bash
# Disable features (prefix with -)
STREMTHRU_FEATURE=-stremio_store,-stremio_wrap,-stremio_sidekick

# Enable specific features only
STREMTHRU_FEATURE=stremio_store,stremio_wrap
```

| Feature | Description | Default |
|---------|-------------|---------|
| `stremio_store` | Stremio Store addon | Enabled |
| `stremio_wrap` | Stremio Wrap addon | Enabled |
| `stremio_sidekick` | Stremio Sidekick addon | Enabled |
| `stremio_list` | Stremio List addon | Enabled |
| `stremio_torz` | Stremio Torz addon | Enabled |
| `imdb_title` | IMDB title syncing | Enabled |
| `dmm_hashlist` | DMM hashlist syncing | Enabled |
| `vault` | Vault features | Enabled |

**Proxy-Only Mode Configuration:**
```bash
# Disable all non-essential features for minimal proxy
STREMTHRU_FEATURE=-stremio_list,-stremio_sidekick,-stremio_store,-stremio_torz,-stremio_wrap,-imdb_title,-dmm_hashlist,-vault
```

---

## Cloudflare WARP Integration

### WARP Overview

Cloudflare WARP integration routes all outbound traffic through Cloudflare's network, ensuring debrid services see a Cloudflare IP instead of your server's IP.

```mermaid
sequenceDiagram
    participant User as Stremio User
    participant ST as StremThru
    participant WP as WireProxy<br/>(SOCKS5)
    participant CF as Cloudflare<br/>WARP
    participant RD as RealDebrid

    User->>ST: Request stream
    ST->>WP: Forward via SOCKS5
    WP->>CF: WireGuard tunnel
    CF->>RD: API request (CF IP: 104.28.x.x)
    RD-->>CF: Response
    CF-->>WP: Response
    WP-->>ST: Response
    ST-->>User: Stream content

    Note over RD: RealDebrid only sees<br/>Cloudflare IP
```

### WARP Setup Guide

#### Step 1: Generate WARP Credentials

```bash
# Download wgcf
curl -Lo wgcf https://github.com/ViRb3/wgcf/releases/latest/download/wgcf_2.2.22_darwin_amd64
chmod +x wgcf

# Register WARP account
./wgcf register --accept-tos

# Generate WireGuard profile
./wgcf generate

# View the generated config
cat wgcf-profile.conf
```

#### Step 2: Create WireProxy Configuration

Create `wireproxy.conf`:

```toml
[Interface]
PrivateKey = YOUR_PRIVATE_KEY_FROM_WGCF
Address = 172.16.0.2/32
Address = 2606:4700:110:xxxx::/128
DNS = 1.1.1.1

[Peer]
PublicKey = bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = engage.cloudflareclient.com:2408

[Socks5]
BindAddress = 127.0.0.1:1080
```

#### Step 3: Create Entrypoint Script

Create `warp-entrypoint.sh`:

```bash
#!/bin/sh
set -e

echo "=== Starting WireProxy (Cloudflare WARP) ==="

# Start wireproxy in background
/app/wireproxy -c /app/wireproxy.conf &

# Wait for SOCKS5 proxy
for i in $(seq 1 30); do
    if nc -z 127.0.0.1 1080 2>/dev/null; then
        echo "SOCKS5 proxy ready on 127.0.0.1:1080"
        break
    fi
    sleep 1
done

# Configure StremThru to use proxy
export STREMTHRU_HTTP_PROXY="socks5://127.0.0.1:1080"

echo "=== Starting StremThru ==="
exec /app/stremthru
```

#### Step 4: Update Dockerfile

Add WireProxy to your Dockerfile:

```dockerfile
# Download WireProxy
FROM alpine AS wireproxy-downloader
RUN apk add --no-cache curl
RUN curl -Lo /wireproxy.tar.gz \
    https://github.com/pufferffish/wireproxy/releases/download/v1.0.9/wireproxy_linux_amd64.tar.gz \
    && tar -xzf /wireproxy.tar.gz -C / \
    && chmod +x /wireproxy

# Final image
FROM alpine
# ... other steps ...
COPY --from=wireproxy-downloader /wireproxy ./wireproxy
COPY wireproxy.conf ./wireproxy.conf
COPY warp-entrypoint.sh ./warp-entrypoint.sh
RUN chmod +x ./warp-entrypoint.sh

CMD ["./warp-entrypoint.sh"]
```

### WARP Configuration

| File | Purpose |
|------|---------|
| `wireproxy.conf` | WireProxy/WARP configuration |
| `warp-entrypoint.sh` | Startup script for WARP + StremThru |
| `wgcf-account.toml` | WARP account credentials (do not commit) |
| `wgcf-profile.conf` | Generated WireGuard profile (do not commit) |

### WARP Troubleshooting

#### WARP Not Connecting

```bash
# Check logs for handshake
heroku logs --tail | grep -i "handshake"

# Expected output:
# peer(bmXO…fgyo) - Sending handshake initiation
# peer(bmXO…fgyo) - Received handshake response
```

#### Verify IP Changed

```bash
# Check exposed IP
curl -H "X-StremThru-Authorization: Basic $(echo -n 'admin:pass' | base64)" \
  https://your-app.herokuapp.com/v0/health/__debug__ | jq '.data.ip.exposed'

# Should show Cloudflare IP (104.28.x.x), not Heroku IP
```

---

## API Reference

### Health Endpoints

#### `GET /v0/health`

Basic health check.

**Response:**
```json
{
  "data": {
    "status": "ok"
  }
}
```

#### `GET /v0/health/__debug__`

Detailed debug information (requires auth).

**Request:**
```bash
curl -H "X-StremThru-Authorization: Basic $(echo -n 'user:pass' | base64)" \
  https://your-instance.com/v0/health/__debug__
```

**Response:**
```json
{
  "data": {
    "time": "2026-02-05T00:00:00Z",
    "version": "0.96.5",
    "user": {
      "name": "admin",
      "store": {
        "default": "realdebrid",
        "names": ["realdebrid"]
      }
    },
    "ip": {
      "client": "YOUR_IP",
      "exposed": {
        "*": "104.28.228.78",
        "download.real-debrid.com": "104.28.228.78"
      },
      "machine": "100.31.199.166",
      "tunnel": {
        "127.0.0.1:1080": "104.28.228.78"
      }
    }
  }
}
```

### Proxy Endpoints

#### `POST /v0/proxy`

Generate proxy URLs for direct links.

**Request:**
```json
{
  "links": [
    "https://download.real-debrid.com/d/XXXXX/file.mkv"
  ]
}
```

**Response:**
```json
{
  "data": {
    "links": [
      "https://your-instance.com/v0/proxy/TOKEN/file.mkv"
    ]
  }
}
```

#### `GET /v0/proxy/{token}`

Access proxied content. Supports byte-range requests for video streaming.

**Headers:**
- `Range`: Standard HTTP range header for byte-serving

### Store Endpoints

All store endpoints require authentication via `X-StremThru-Authorization` header.

#### `GET /v0/store/user`

Get current user info from debrid service.

**Response:**
```json
{
  "data": {
    "id": "13059080",
    "email": "user@example.com",
    "subscription_status": "premium"
  }
}
```

#### `POST /v0/store/magnets`

Add a magnet to the debrid service.

**Request:**
```json
{
  "magnet": "magnet:?xt=urn:btih:HASH..."
}
```

**Response:**
```json
{
  "data": {
    "id": "TORRENT_ID",
    "hash": "HASH",
    "status": "downloading",
    "files": [...]
  }
}
```

#### `GET /v0/store/magnets`

List all magnets in your debrid account.

#### `GET /v0/store/magnets/check`

Check if magnets are cached.

**Query Parameters:**
- `magnets[]`: Array of magnet URIs or hashes

#### `POST /v0/store/link/generate`

Generate direct download link from a file link.

**Request:**
```json
{
  "link": "https://real-debrid.com/file/XXXXX"
}
```

---

## System Design

### Architecture Diagram

```mermaid
flowchart TB
    subgraph Client["Client Layer"]
        Stremio[Stremio App]
        Browser[Web Browser]
    end

    subgraph StremThru["StremThru Application"]
        direction TB
        HTTP[HTTP Server<br/>:8080]
        
        subgraph Endpoints["Endpoint Handlers"]
            Health[/health]
            Proxy[/proxy]
            Store[/store]
            Addons[/stremio/*]
        end
        
        subgraph Middleware["Middleware"]
            Auth[Auth Middleware]
            CORS[CORS Middleware]
            Context[Store Context]
        end
        
        subgraph Core["Core Services"]
            StoreClient[Store Clients]
            ProxyHandler[Proxy Handler]
            Cache[Cache Layer]
        end
        
        subgraph Network["Network Layer"]
            WireProxy[WireProxy<br/>SOCKS5 :1080]
            HTTPClient[HTTP Client]
        end
    end

    subgraph Storage["Storage"]
        SQLite[(SQLite/Postgres)]
        Redis[(Redis<br/>Optional)]
    end

    subgraph External["External Services"]
        WARP[Cloudflare WARP]
        RealDebrid[RealDebrid]
        AllDebrid[AllDebrid]
        Other[Other Debrids...]
    end

    Stremio --> HTTP
    Browser --> HTTP
    
    HTTP --> Middleware
    Middleware --> Endpoints
    Endpoints --> Core
    
    Core --> SQLite
    Core --> Redis
    Core --> Network
    
    Network --> WireProxy
    WireProxy --> WARP
    WARP --> RealDebrid
    WARP --> AllDebrid
    WARP --> Other
```

### Request Flow

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant S as StremThru
    participant W as WireProxy
    participant CF as Cloudflare WARP
    participant RD as RealDebrid

    C->>S: POST /v0/store/magnets<br/>Authorization: Basic xxx
    
    activate S
    S->>S: Validate Auth
    S->>S: Parse Request
    S->>S: Get Store Context
    
    S->>W: SOCKS5 Connect
    activate W
    W->>CF: WireGuard Tunnel
    activate CF
    CF->>RD: POST /torrents/addMagnet<br/>IP: 104.28.x.x
    activate RD
    RD-->>CF: Torrent Added
    deactivate RD
    CF-->>W: Response
    deactivate CF
    W-->>S: Response
    deactivate W
    
    S-->>C: {"data": {"id": "xxx", "status": "downloading"}}
    deactivate S
```

### Data Flow

```mermaid
flowchart LR
    subgraph Input
        Magnet[Magnet URI]
        Torrent[.torrent File]
        Link[Direct Link]
    end
    
    subgraph Process
        Parse[Parse Input]
        Check[Check Cache]
        Add[Add to Debrid]
        Generate[Generate Link]
        Proxify[Create Proxy URL]
    end
    
    subgraph Output
        ProxyURL[Proxy URL]
        Stream[Video Stream]
    end
    
    Magnet --> Parse
    Torrent --> Parse
    Link --> Parse
    
    Parse --> Check
    Check -->|Cached| Generate
    Check -->|Not Cached| Add
    Add --> Generate
    Generate --> Proxify
    Proxify --> ProxyURL
    ProxyURL --> Stream
```

### Component Interactions

| Component | Responsibility | Interacts With |
|-----------|---------------|----------------|
| `main.go` | Application entry point, server setup | All components |
| `internal/endpoint/*` | HTTP handlers for API endpoints | Middleware, Store clients |
| `internal/config/*` | Configuration management | All components |
| `internal/db/*` | Database operations | SQLite/PostgreSQL |
| `store/*` | Debrid service clients | External APIs |
| `internal/shared/*` | Shared utilities, proxy handler | HTTP client |
| `WireProxy` | SOCKS5 proxy via WARP | Cloudflare network |

---

## Commands Reference

### Go Commands

```bash
# Run the application
go run --tags "fts5" .

# Build binary
go build --tags "fts5" -o stremthru

# Run tests
STREMTHRU_ENV=test go test -v ./...

# Format code
go fmt ./...

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### Make Commands

```bash
# Build everything
make all

# Build binary only
make build

# Run locally
make run

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Push Docker image
make docker-push
```

### Node/PNPM Commands

```bash
# Install dependencies
pnpm install

# Build dashboard
pnpm run dash:build

# Run dashboard dev server
pnpm run dash:dev

# Lint staged files
pnpm run lint:staged

# Lint current commit
pnpm run lint:commit:current
```

### Heroku Commands

```bash
# Login to Heroku
heroku login

# Create app
heroku create your-app-name --stack container

# Set config vars
heroku config:set KEY=VALUE -a your-app-name

# View config
heroku config -a your-app-name

# Deploy
git push heroku main

# View logs
heroku logs --tail -a your-app-name

# Restart app
heroku restart -a your-app-name

# Open app
heroku open -a your-app-name

# Check dyno status
heroku ps -a your-app-name

# Run one-off command
heroku run bash -a your-app-name
```

---

## Operations

### Logging

StremThru uses structured JSON logging:

```json
{"time":"2026-02-05T00:00:00Z","level":"INFO","msg":"stremthru listening on :8080"}
```

**Log Levels:**
- `DEBUG`: Detailed debugging information
- `INFO`: General operational messages
- `WARN`: Warning messages
- `ERROR`: Error messages

**Viewing Logs:**
```bash
# Docker
docker logs -f stremthru

# Heroku
heroku logs --tail -a your-app-name

# Local
go run . 2>&1 | jq .
```

### Monitoring

#### Health Check

```bash
# Basic health
curl https://your-instance.com/v0/health

# Detailed debug (with auth)
curl -H "X-StremThru-Authorization: Basic xxx" \
  https://your-instance.com/v0/health/__debug__
```

#### Uptime Monitoring

Configure external monitoring services (UptimeRobot, Pingdom) to hit:
- Endpoint: `/v0/health`
- Expected response: `{"data":{"status":"ok"}}`
- Interval: 5 minutes

### Database Management

#### SQLite (Default)

```bash
# Location
./data/stremthru.db

# Backup
cp ./data/stremthru.db ./data/stremthru.db.backup

# Connect with CLI
sqlite3 ./data/stremthru.db
```

#### PostgreSQL

```bash
# Connect
psql $STREMTHRU_DATABASE_URI

# Backup
pg_dump $STREMTHRU_DATABASE_URI > backup.sql

# Restore
psql $STREMTHRU_DATABASE_URI < backup.sql
```

### Scaling Considerations

| Factor | Recommendation |
|--------|---------------|
| **Concurrent Users** | Each user maintains connections to debrid APIs. Monitor API rate limits. |
| **Memory** | ~50-100MB base. Increases with active streams. |
| **CPU** | Low usage. Primarily I/O bound. |
| **Storage** | SQLite file grows with metadata. Consider PostgreSQL for high volume. |
| **Network** | Video content is proxied. Ensure sufficient bandwidth. |

---

## Security

### Authentication

StremThru uses HTTP Basic Authentication:

```bash
# Header format
Authorization: Basic base64(username:password)

# Or StremThru-specific header
X-StremThru-Authorization: Basic base64(username:password)

# Example
curl -H "Authorization: Basic $(echo -n 'admin:password' | base64)" \
  https://your-instance.com/v0/store/user
```

### IP Privacy Model

```mermaid
flowchart LR
    subgraph "What Each Service Sees"
        direction TB
        
        subgraph You["Your Device"]
            YourIP["Your IP: 1.2.3.4"]
        end
        
        subgraph ST["StremThru"]
            STSees["Sees: 1.2.3.4"]
        end
        
        subgraph WARP["Cloudflare WARP"]
            WARPSees["Sees: StremThru Internal"]
        end
        
        subgraph RD["RealDebrid"]
            RDSees["Sees: 104.28.x.x<br/>(Cloudflare IP)"]
        end
    end
    
    You --> ST
    ST --> WARP
    WARP --> RD
    
    style RDSees fill:#c8e6c9
```

**Key Points:**
- Your IP is visible only to StremThru
- RealDebrid never sees your IP
- With WARP, RealDebrid sees Cloudflare's IP
- Without WARP, RealDebrid sees StremThru's server IP

### Secrets Handling

| Secret | Storage | Notes |
|--------|---------|-------|
| `STREMTHRU_PROXY_AUTH` | Environment variable | Never commit to git |
| `STREMTHRU_STORE_AUTH` | Environment variable | Contains debrid API keys |
| `wireproxy.conf` | File (gitignored) | Contains WARP private key |
| `wgcf-account.toml` | File (gitignored) | WARP account credentials |

**Best Practices:**
1. Use environment variables for secrets
2. Never commit `.env` files
3. Use Heroku config vars or Docker secrets in production
4. Rotate API keys periodically

---

## Troubleshooting

### Common Issues

#### Application Won't Start

```bash
# Check logs
docker logs stremthru

# Verify environment variables are set
heroku config -a your-app-name

# Check database connectivity
sqlite3 ./data/stremthru.db "SELECT 1;"
```

#### RealDebrid Connection Failed

```bash
# Verify API token
curl -H "Authorization: Bearer YOUR_RD_TOKEN" \
  https://api.real-debrid.com/rest/1.0/user

# Check StremThru store config
curl -H "X-StremThru-Authorization: Basic xxx" \
  https://your-instance.com/v0/health/__debug__ | jq '.data.user'
```

#### WARP Not Working

```bash
# Check WireProxy logs
heroku logs --tail | grep -i "peer"

# Verify SOCKS5 proxy is running
heroku run "nc -z 127.0.0.1 1080" -a your-app-name

# Check exposed IP
curl -H "X-StremThru-Authorization: Basic xxx" \
  https://your-instance.com/v0/health/__debug__ | jq '.data.ip.exposed'
```

#### Stream Not Playing

```bash
# Test proxy URL directly
curl -I "https://your-instance.com/v0/proxy/TOKEN/file.mkv"

# Check for byte-range support
curl -H "Range: bytes=0-1000" -I "https://your-instance.com/v0/proxy/TOKEN/file.mkv"
```

---

## Contributing

### Repository Structure

```
stremthru/
├── apps/
│   └── dash/              # Frontend dashboard (TanStack Start)
├── core/                  # Core utilities
│   ├── error.go          # Error handling
│   ├── magnet.go         # Magnet link parsing
│   └── request.go        # HTTP utilities
├── internal/
│   ├── config/           # Configuration management
│   ├── db/               # Database layer
│   ├── endpoint/         # HTTP endpoint handlers
│   ├── shared/           # Shared utilities
│   └── worker/           # Background workers
├── migrations/           # Database migrations
├── store/
│   ├── realdebrid/      # RealDebrid client
│   ├── alldebrid/       # AllDebrid client
│   └── .../             # Other debrid clients
├── stremio/             # Stremio addon definitions
├── Dockerfile           # Multi-stage Docker build
├── app.json             # Heroku app manifest
├── heroku.yml           # Heroku container config
├── wireproxy.conf       # WARP proxy config
└── warp-entrypoint.sh   # WARP startup script
```

### Coding Standards

#### Go

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `go fmt` for formatting
- Run `go vet` before committing
- Write table-driven tests

#### TypeScript (Dashboard)

- Use ESLint configuration provided
- Follow Prettier formatting
- Use TypeScript strict mode

#### Commits

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
feat: add support for new debrid service
fix: resolve proxy timeout issue
docs: update README with WARP setup
chore: update dependencies
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes
4. Run tests: `make test`
5. Format code: `make fmt`
6. Commit with conventional commit message
7. Push and create Pull Request
8. Wait for CI checks and review

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Made with ❤️ for the streaming community**

[Report Bug](https://github.com/MunifTanjim/stremthru/issues) · [Request Feature](https://github.com/MunifTanjim/stremthru/issues)

</div>
