package endpoint

const statusHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>StremThru Status Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            color: #eee;
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        h1 {
            font-size: 2rem;
            margin-bottom: 24px;
            background: linear-gradient(90deg, #00d4ff, #7b68ee);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-bottom: 24px; }
        .card {
            background: rgba(255,255,255,0.05);
            border: 1px solid rgba(255,255,255,0.1);
            border-radius: 12px;
            padding: 20px;
            backdrop-filter: blur(10px);
        }
        .card h2 { font-size: 0.9rem; text-transform: uppercase; letter-spacing: 1px; color: #888; margin-bottom: 12px; }
        .card h3 { font-size: 1.5rem; margin-bottom: 8px; }
        .status-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid rgba(255,255,255,0.05); }
        .status-row:last-child { border-bottom: none; }
        .status-label { color: #888; }
        .status-value { font-family: 'Monaco', monospace; color: #00d4ff; }
        .status-ok { color: #4ade80; }
        .status-error { color: #f87171; }
        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.8rem;
            font-weight: 600;
        }
        .badge-ok { background: rgba(74,222,128,0.2); color: #4ade80; }
        .badge-error { background: rgba(248,113,113,0.2); color: #f87171; }
        .badge-warn { background: rgba(251,191,36,0.2); color: #fbbf24; }
        .btn {
            background: linear-gradient(90deg, #00d4ff, #7b68ee);
            border: none;
            padding: 12px 24px;
            border-radius: 8px;
            color: #fff;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s, opacity 0.2s;
        }
        .btn:hover { transform: translateY(-2px); }
        .btn:disabled { opacity: 0.5; cursor: not-allowed; transform: none; }
        .speed-result { margin-top: 16px; }
        .speed-value { font-size: 2.5rem; font-weight: 700; }
        .speed-unit { font-size: 1rem; color: #888; }
        table { width: 100%; border-collapse: collapse; font-size: 0.85rem; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid rgba(255,255,255,0.1); }
        th { color: #888; font-weight: 500; text-transform: uppercase; font-size: 0.75rem; }
        .loading { color: #888; font-style: italic; }
        .error-box { background: rgba(248,113,113,0.1); border: 1px solid #f87171; padding: 12px; border-radius: 8px; margin-top: 12px; }
        .segment { margin: 12px 0; padding: 12px; background: rgba(0,0,0,0.2); border-radius: 8px; }
        .segment-label { font-size: 0.8rem; color: #888; margin-bottom: 8px; }
        .segment-speed { font-size: 1.5rem; font-weight: 600; color: #00d4ff; }
        .ip-chain { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; margin: 8px 0; }
        .ip-node { background: rgba(0,212,255,0.1); border: 1px solid rgba(0,212,255,0.3); padding: 8px 12px; border-radius: 8px; text-align: center; }
        .ip-node-label { font-size: 0.7rem; color: #888; }
        .ip-node-value { font-family: monospace; font-size: 0.9rem; color: #00d4ff; }
        .ip-arrow { color: #888; font-size: 1.2rem; }
        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
        .pulsing { animation: pulse 1.5s infinite; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌐 StremThru Network Status</h1>
        
        <!-- IP Chain -->
        <div class="card" style="margin-bottom: 20px;">
            <h2>Network Path</h2>
            <div class="ip-chain" id="ipChain">
                <div class="ip-node">
                    <div class="ip-node-label">Your IP</div>
                    <div class="ip-node-value" id="clientIp">--</div>
                </div>
                <div class="ip-arrow">→</div>
                <div class="ip-node">
                    <div class="ip-node-label">Server (Heroku)</div>
                    <div class="ip-node-value" id="machineIp">--</div>
                </div>
                <div class="ip-arrow">→</div>
                <div class="ip-node">
                    <div class="ip-node-label">WARP/Cloudflare</div>
                    <div class="ip-node-value" id="tunnelIp">--</div>
                    <div class="ip-node-label" id="tunnelAsn"></div>
                </div>
                <div class="ip-arrow">→</div>
                <div class="ip-node">
                    <div class="ip-node-label">RealDebrid sees</div>
                    <div class="ip-node-value" id="rdSeenIp">--</div>
                </div>
            </div>
        </div>

        <div class="grid">
            <!-- Status Card -->
            <div class="card">
                <h2>Connection Status</h2>
                <div class="status-row">
                    <span class="status-label">WARP Active</span>
                    <span id="warpStatus" class="badge badge-warn">Checking...</span>
                </div>
                <div class="status-row">
                    <span class="status-label">RealDebrid</span>
                    <span id="rdStatus" class="badge badge-warn">Checking...</span>
                </div>
                <div class="status-row">
                    <span class="status-label">RD Latency</span>
                    <span class="status-value" id="rdLatency">--</span>
                </div>
                <div class="status-row">
                    <span class="status-label">Last Check</span>
                    <span class="status-value" id="lastCheck">--</span>
                </div>
                <div id="statusError" class="error-box" style="display:none;"></div>
            </div>

            <!-- Speed Test Card -->
            <div class="card">
                <h2>Speed Test</h2>
                <button class="btn" id="speedTestBtn" onclick="runSpeedTest()">Run Speed Test</button>
                <div class="speed-result" id="speedResult" style="display:none;">
                    <div class="segment">
                        <div class="segment-label">Server → WARP (Cloudflare)</div>
                        <div class="segment-speed" id="warpSpeed">--</div>
                        <div style="color:#888;font-size:0.8rem;">Latency: <span id="warpLatency">--</span>ms</div>
                    </div>
                    <div class="segment">
                        <div class="segment-label">WARP → RealDebrid</div>
                        <div class="segment-speed" id="rdSpeed">--</div>
                        <div style="color:#888;font-size:0.8rem;">Latency: <span id="rdTestLatency">--</span>ms</div>
                    </div>
                </div>
                <div id="speedError" class="error-box" style="display:none;"></div>
            </div>

            <!-- Bandwidth Card -->
            <div class="card">
                <h2>Data Transferred</h2>
                <div class="status-row">
                    <span class="status-label">Dashboard In</span>
                    <span class="status-value" id="bytesIn">--</span>
                </div>
                <div class="status-row">
                    <span class="status-label">Dashboard Out</span>
                    <span class="status-value" id="bytesOut">--</span>
                </div>
                <div class="status-row">
                    <span class="status-label">Speed Tests</span>
                    <span class="status-value" id="speedBytes">--</span>
                </div>
                <div class="status-row">
                    <span class="status-label">Since Restart</span>
                    <span class="status-value" id="uptime">--</span>
                </div>
            </div>

            <!-- Client Latency Card -->
            <div class="card">
                <h2>Your Connection</h2>
                <div class="status-row">
                    <span class="status-label">Client → Server</span>
                    <span class="status-value" id="clientLatency">--</span>
                </div>
                <button class="btn" style="margin-top:12px;background:#444;" onclick="measureClientLatency()">Measure Latency</button>
            </div>
        </div>

        <!-- History -->
        <div class="card">
            <h2>Recent Speed Tests</h2>
            <table>
                <thead>
                    <tr>
                        <th>Time</th>
                        <th>Server→WARP</th>
                        <th>WARP→RD</th>
                        <th>Total Latency</th>
                    </tr>
                </thead>
                <tbody id="historyTable">
                    <tr><td colspan="4" class="loading">Loading...</td></tr>
                </tbody>
            </table>
        </div>
    </div>

    <script>
        function formatBytes(bytes) {
            if (!bytes || bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        function formatTime(dateStr) {
            if (!dateStr) return '--';
            const d = new Date(dateStr);
            return d.toLocaleTimeString();
        }

        function formatDuration(start) {
            const ms = Date.now() - new Date(start).getTime();
            const hours = Math.floor(ms / 3600000);
            const mins = Math.floor((ms % 3600000) / 60000);
            return hours + 'h ' + mins + 'm';
        }

        async function fetchStatus() {
            try {
                const res = await fetch('/status/api/current');
                const data = await res.json();
                
                document.getElementById('machineIp').textContent = data.machine_ip || '--';
                document.getElementById('tunnelIp').textContent = data.tunnel_ip || '--';
                document.getElementById('tunnelAsn').textContent = data.tunnel_asn || '';
                document.getElementById('rdSeenIp').textContent = data.realdebrid_seen_ip || '--';
                
                const warpEl = document.getElementById('warpStatus');
                warpEl.textContent = data.warp_active ? 'Active' : 'Inactive';
                warpEl.className = 'badge ' + (data.warp_active ? 'badge-ok' : 'badge-error');
                
                const rdEl = document.getElementById('rdStatus');
                rdEl.textContent = data.realdebrid_ok ? 'Connected' : 'Error';
                rdEl.className = 'badge ' + (data.realdebrid_ok ? 'badge-ok' : 'badge-error');
                
                document.getElementById('rdLatency').textContent = data.realdebrid_latency_ms + ' ms';
                document.getElementById('lastCheck').textContent = formatTime(data.checked_at);
                
                if (data.last_error) {
                    document.getElementById('statusError').textContent = data.last_error;
                    document.getElementById('statusError').style.display = 'block';
                } else {
                    document.getElementById('statusError').style.display = 'none';
                }
            } catch (e) {
                console.error('Status fetch error:', e);
            }
        }

        async function fetchBandwidth() {
            try {
                const res = await fetch('/status/api/bandwidth');
                const data = await res.json();
                document.getElementById('bytesIn').textContent = formatBytes(data.dashboard_bytes_in);
                document.getElementById('bytesOut').textContent = formatBytes(data.dashboard_bytes_out);
                document.getElementById('speedBytes').textContent = formatBytes(data.speedtest_bytes);
                document.getElementById('uptime').textContent = formatDuration(data.since_restart_at);
            } catch (e) {
                console.error('Bandwidth fetch error:', e);
            }
        }

        async function fetchHistory() {
            try {
                const res = await fetch('/status/api/history');
                const data = await res.json();
                const tbody = document.getElementById('historyTable');
                
                if (!data.speed_tests || data.speed_tests.length === 0) {
                    tbody.innerHTML = '<tr><td colspan="4" style="color:#888;">No speed tests yet</td></tr>';
                    return;
                }
                
                tbody.innerHTML = data.speed_tests.slice(0, 10).map(t => {
                    const warp = t.server_to_warp ? t.server_to_warp.speed_mbps.toFixed(1) + ' Mbps' : '--';
                    const rd = t.warp_to_realdebrid ? t.warp_to_realdebrid.speed_mbps.toFixed(1) + ' Mbps' : '--';
                    return '<tr><td>' + formatTime(t.tested_at) + '</td><td>' + warp + '</td><td>' + rd + '</td><td>' + t.total_latency_ms + ' ms</td></tr>';
                }).join('');
            } catch (e) {
                console.error('History fetch error:', e);
            }
        }

        async function runSpeedTest() {
            const btn = document.getElementById('speedTestBtn');
            btn.disabled = true;
            btn.textContent = 'Testing...';
            btn.classList.add('pulsing');
            
            document.getElementById('speedError').style.display = 'none';
            
            try {
                const res = await fetch('/status/api/speedtest', { method: 'POST' });
                const data = await res.json();
                
                document.getElementById('speedResult').style.display = 'block';
                
                if (data.server_to_warp) {
                    document.getElementById('warpSpeed').textContent = data.server_to_warp.speed_mbps.toFixed(1) + ' Mbps';
                    document.getElementById('warpLatency').textContent = data.server_to_warp.latency_ms;
                }
                
                if (data.warp_to_realdebrid) {
                    document.getElementById('rdSpeed').textContent = data.warp_to_realdebrid.speed_mbps.toFixed(1) + ' Mbps';
                    document.getElementById('rdTestLatency').textContent = data.warp_to_realdebrid.latency_ms;
                }
                
                if (data.error) {
                    document.getElementById('speedError').textContent = data.error;
                    document.getElementById('speedError').style.display = 'block';
                }
                
                fetchHistory();
                fetchBandwidth();
            } catch (e) {
                document.getElementById('speedError').textContent = 'Speed test failed: ' + e.message;
                document.getElementById('speedError').style.display = 'block';
            } finally {
                btn.disabled = false;
                btn.textContent = 'Run Speed Test';
                btn.classList.remove('pulsing');
            }
        }

        async function measureClientLatency() {
            const samples = [];
            for (let i = 0; i < 5; i++) {
                const start = performance.now();
                try {
                    await fetch('/status/api/ping');
                    samples.push(performance.now() - start);
                } catch (e) {}
            }
            if (samples.length > 0) {
                const avg = samples.reduce((a, b) => a + b) / samples.length;
                document.getElementById('clientLatency').textContent = avg.toFixed(0) + ' ms';
            }
        }

        // Detect client IP via fetch headers (server will see it)
        document.getElementById('clientIp').textContent = 'Your Browser';

        // Initial load
        fetchStatus();
        fetchBandwidth();
        fetchHistory();
        measureClientLatency();

        // Auto-refresh every 30s
        setInterval(fetchStatus, 30000);
        setInterval(fetchBandwidth, 30000);
    </script>
</body>
</html>`
