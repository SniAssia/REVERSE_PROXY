# REVERSE_PROXY
# Reverse Proxy Load Balancer

A mini infrastructure system written in **Go** that implements a reverse proxy with:

- Load balancing (Round‑Robin & Least‑Connections)
- Backend health checking
- Admin HTTP API for dynamic control
- Graceful shutdown
- Concurrent request handling

This is a real-world infrastructure component — not a toy project.

---

## 🧠 Concepts

- **Proxy Server** listens on `:8080` and forwards incoming HTTP traffic to backends.
- **Load Balancer** selects a backend using either round‑robin or least‑connections.
- **Health Checker** monitors backends via their `/health` endpoint.
- **Admin API** on `:9000` allows dynamic management of backends and strategy.
- **Graceful Shutdown** handles SIGINT/SIGTERM for clean exit.

---

##  Getting Started

### Requirements

- Go 1.18+
- `curl` (for testing)

---

##  Steps to Run

### 1️⃣ Mock Backends

In three separate terminals run:

```bash
go run cmd/backend/main.go 8081
go run cmd/backend/main.go 8082
go run cmd/backend/main.go 8083

### 2️⃣ Start Proxy + Admin

From the cmd/proxy directory:

go run main.go handler.go


Output:

[proxy] listening on :8080
[admin] listening on :9000

⚙️ Admin API
➕ Add Backend
curl -X POST http://localhost:9000/backends \
     -H "Content-Type: application/json" \
     -d '{"url":"http://localhost:8081"}'


Add all backends one by one.

📊 Get Status
curl http://localhost:9000/status

Response includes:

current strategy (round-robin / least-connections)

each backend’s alive status and active connections

##EXAMPLE 
{
  "strategy":"round-robin",
  "backends":[
    {"url":"http://localhost:8081","alive":true,"conns":0},
    {"url":"http://localhost:8082","alive":true,"conns":0}
  ]
}
🔀 Switch Strategy
 // still working on it 
Round‑Robin
curl -X POST http://localhost:9000/strategy \
     -H "Content-Type: application/json" \
     -d '{"strategy":"round-robin"}'

Least‑Connections
curl -X POST http://localhost:9000/strategy \
     -H "Content-Type: application/json" \
     -d '{"strategy":"least-connections"}'

➖ Remove Backend
curl -X DELETE http://localhost:9000/backends \
     -H "Content-Type: application/json" \
     -d '{"url":"http://localhost:8082"}'

 Proxy Usage
curl http://localhost:8080/

🛠 Error Handling & Edge Cases
Invalid JSON
curl -X POST http://localhost:9000/backends \
     -H "Content-Type: application/json" \
     -d '{bad json}'

