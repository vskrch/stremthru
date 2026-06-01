#!/bin/bash
set -e

echo "=== StremThru + WARP Startup ==="

# Create tun device if not exists
if [ ! -e /dev/net/tun ]; then
    sudo mkdir -p /dev/net
    sudo mknod /dev/net/tun c 10 200
    sudo chmod 600 /dev/net/tun
fi

# Start dbus
sudo mkdir -p /run/dbus
if [ -f /run/dbus/pid ]; then
    sudo rm /run/dbus/pid
fi
sudo dbus-daemon --config-file=/usr/share/dbus-1/system.conf 2>/dev/null || true

# Start WARP daemon
echo "Starting WARP daemon..."
sudo warp-svc --accept-tos &

# Wait for daemon
sleep "${WARP_SLEEP:-5}"

# Register if needed
if [ ! -f /var/lib/cloudflare-warp/reg.json ]; then
    if [ ! -f /var/lib/cloudflare-warp/mdm.xml ] || [ -n "${REGISTER_WHEN_MDM_EXISTS}" ]; then
        echo "Registering WARP client..."
        warp-cli registration new && echo "WARP client registered!"
        if [ -n "${WARP_LICENSE_KEY}" ]; then
            echo "Registering WARP license..."
            warp-cli registration license "${WARP_LICENSE_KEY}" && echo "WARP license registered!"
        fi
    fi
fi

# Set WARP to proxy mode (no tun device needed)
echo "Configuring WARP proxy mode..."
warp-cli --accept-tos mode proxy
warp-cli --accept-tos proxy port 40000

# Connect
echo "Connecting to WARP..."
warp-cli --accept-tos connect

# Disable qlog
warp-cli --accept-tos debug qlog disable

# Wait for WARP to connect
sleep 3

# Verify WARP connection
echo "Verifying WARP connection..."
WARP_STATUS=$(curl -fsS --socks5-hostname 127.0.0.1:40000 "https://cloudflare.com/cdn-cgi/trace" 2>/dev/null | grep "warp=" || echo "warp=unknown")
echo "WARP status: ${WARP_STATUS}"

# Start gost proxy chain
echo "Starting gost proxy..."
gost ${GOST_ARGS} &
GOST_PID=$!

# Wait for gost to start
sleep 2

# Verify gost is working
echo "Verifying gost proxy..."
if curl -fsS --socks5-hostname 127.0.0.1:1080 "https://cloudflare.com/cdn-cgi/trace" > /dev/null 2>&1; then
    echo "gost proxy verified OK"
else
    echo "WARNING: gost proxy verification failed, starting anyway..."
fi

# Set proxy for stremthru
export STREMTHRU_HTTP_PROXY="socks5://127.0.0.1:1080"
export STREMTHRU_TUNNEL="*:true"

echo "Starting StremThru on port ${PORT:-8080}..."
exec stremthru
