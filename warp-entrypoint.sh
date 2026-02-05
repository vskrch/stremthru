#!/bin/sh
set -e

echo "=== Starting WireProxy (Cloudflare WARP) ==="

# Check if wireproxy config exists
if [ ! -f "/app/wireproxy.conf" ]; then
    echo "ERROR: wireproxy.conf not found!"
    exit 1
fi

# Start wireproxy in background
/app/wireproxy -c /app/wireproxy.conf &
WIREPROXY_PID=$!

# Wait for SOCKS5 proxy to be ready
echo "Waiting for SOCKS5 proxy to be ready..."
for i in $(seq 1 30); do
    if nc -z 127.0.0.1 1080 2>/dev/null; then
        echo "SOCKS5 proxy is ready on 127.0.0.1:1080"
        break
    fi
    sleep 1
done

# Verify proxy is running
if ! nc -z 127.0.0.1 1080 2>/dev/null; then
    echo "ERROR: SOCKS5 proxy failed to start!"
    exit 1
fi

# Export proxy for StremThru
export STREMTHRU_HTTP_PROXY="socks5://127.0.0.1:1080"
echo "STREMTHRU_HTTP_PROXY set to: $STREMTHRU_HTTP_PROXY"

echo "=== Starting StremThru ==="
exec /app/stremthru
