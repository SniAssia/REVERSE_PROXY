1️⃣ Project Architecture – Components Overview

Your system is composed of three logical services plus shared core logic.

A. Reverse Proxy Server (Public Entry Point)

Purpose:
Accept client traffic and forward it to backend services using a load-balancing strategy.

Responsibilities:

Accept HTTP requests from clients

Select a healthy backend

Forward requests using httputil.ReverseProxy

Track active connections

Handle backend failures

Propagate request contexts

Listens on: :8080

B. Server Pool & Load Balancer (Core Logic)

Purpose:
Maintain backend state and decide where traffic goes.

Responsibilities:

Store backend servers

Track health and connection counts

Implement Round-Robin / Least-Connections

Ensure thread safety

Expose backend selection via interface

C. Health Checker (Background Worker)

Purpose:
Continuously verify backend availability.

Responsibilities:

Periodically ping /health

Update backend Alive status

Log transitions (UP → DOWN)

Integrate with ServerPool safely

D. Admin API (Control Plane)

Purpose:
Dynamically manage the system at runtime.

Responsibilities:

Add / remove backends

View backend status

Change load-balancing strategy

Expose metrics (optional)

Listens on: :8081

E. Backend APIs (Mock Services)

Purpose:
Simulate real services behind the proxy.

Responsibilities:

Respond to proxied requests

Provide /health endpoint

Run on separate ports

Ports: :8082, :8083, :8084, etc.

2️⃣ System Architecture Diagram (Textual)
                    ┌─────────────────────┐
                    │        Client        │
                    └─────────┬───────────┘
                              │
                              │ HTTP Requests
                              ▼
               ┌────────────────────────────────┐
               │       Reverse Proxy (:8080)     │
               │--------------------------------│
               │  • Request Handler              │
               │  • Load Balancer Interface      │
               │  • ReverseProxy                 │
               └───────────┬────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Backend A   │    │ Backend B   │    │ Backend C   │
│ :8082       │    │ :8083       │    │ :8084       │
│ /health     │    │ /health     │    │ /health     │
└─────────────┘    └─────────────┘    └─────────────┘

     ▲
     │
┌────┴─────────────────────────────────────────┐
│          Health Checker (Goroutine)           │
│  • Periodic health pings                      │
│  • Updates backend state                      │
└──────────────────────────────────────────────┘

┌──────────────────────────────────────────────┐
│           Admin API (:8081)                   │
│  • GET /status                                │
│  • POST /backends                             │
│  • DELETE /backends                           │
│  • POST /strategy                             │
└──────────────────────────────────────────────┘
