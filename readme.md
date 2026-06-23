# GoCache Utility Over TCP

A lightweight Redis-inspired key-value cache built in Go with TCP networking support.

GoCache Utility Over TCP is a client-server cache system that enables applications to store and retrieve key-value pairs over a TCP connection. The server maintains an in-memory cache for fast access while periodically persisting data to disk and automatically compacting storage in the background.

This project was built to explore concepts behind distributed cache systems such as Redis, including networking, persistence, startup recovery, concurrent client handling, synchronization, and storage optimization.

---

## Features

* TCP-based client-server architecture
* In-memory key-value storage
* CSV-backed persistence
* Automatic periodic snapshots
* Automatic scheduled compaction
* Startup recovery from persisted data
* Graceful shutdown handling
* JSON-based request protocol
* Concurrent client handling using Goroutines
* Thread-safe operations using RWMutex
* Fast O(1) average lookups using Go maps
* Simple client constructor API

---

## Architecture

```text
+------------+
|   Client   |
+------------+
       |
       | TCP
       v
+------------------+
|   Cache Server   |
+------------------+
       |
       v
+------------------+
| In-Memory Cache  |
| map[string]string|
+------------------+
       |
       v
+------------------+
|     db.csv       |
+------------------+
```

The server keeps all active records in memory while periodically persisting them to disk.

Clients communicate with the server using JSON messages over TCP connections.

---

## How It Works

### In-Memory Storage

The cache is maintained using:

```go
map[string]string
```

This provides near O(1) average lookup and insertion performance.

---

### Thread Safety

Since multiple clients may access the cache simultaneously, synchronization is achieved using:

```go
var mut sync.RWMutex
```

This ensures safe concurrent reads and writes.

---

### TCP Communication

Clients send JSON requests to the server.

Example request:

```json
{
  "cmd": "SET",
  "key": "username",
  "value": "Ankit"
}
```

The server processes the request and sends a response over the same connection.

---

### Startup Recovery

When the server starts, it loads existing data from `db.csv` and reconstructs the in-memory cache.

```go
LoadData()
```

This ensures data survives server restarts.

---

### Automatic Persistence

The server periodically writes the current cache state to disk.

Current configuration:

```go
tickerWriteDb := time.NewTicker(5 * time.Second)
```

Benefits:

* Reduces risk of data loss
* No manual save operation required
* Automatic persistence in the background

---

### Automatic Compaction

The server automatically compacts the database file.

Current configuration:

```go
tickerCompact := time.NewTicker(11 * time.Second)
```

Compaction rewrites storage using only the latest values currently held in memory.

Benefits:

* Smaller storage size
* Faster recovery times
* Reduced duplicate entries

---

### Graceful Shutdown

The server listens for:

```text
SIGINT
SIGTERM
```

Before shutting down it:

1. Persists the latest cache state
2. Performs a final compaction
3. Exits safely

This helps prevent data loss.

---

## Supported Commands

### SET

Store a key-value pair.

Request:

```json
{
  "cmd": "SET",
  "key": "name",
  "value": "YourName"
}
```

Response:

```text
OK
```

---

### GET

Retrieve a value by key.

Request:

```json
{
  "cmd": "GET",
  "key": "name"
}
```

Response:

```text
YourName
```

---

## Installation

Install the package:

```bash
go get github.com/AnkitDTUSE/GocacheUtilityOverTcp@latest
```

Update dependencies:

```bash
go mod tidy
```

---

## Library Structure

```text
GocacheUtilityOverTcp
├── client
│   └── client.go
│
├── server
│   ├── cacheUtil.go
│   ├── compact.go
│   └── server.go
│
├── go.mod
├── go.sum
└── README.md
```

---

## Starting the Server

Create:

```go
package main

import (
	"fmt"

	s "github.com/AnkitDTUSE/GocacheUtilityOverTcp/server"
)

func main() {

	err := s.Start(3000, "tcp")

	if err != nil {
		fmt.Println("error while starting server")
	}
}
```

Run:

```bash
go run .
```

The server will:

* Start listening on port 3000
* Load existing data from disk
* Accept multiple client connections
* Automatically persist data
* Automatically compact storage

---

## Creating a Client

```go
package main

import (
	"fmt"

	c "github.com/AnkitDTUSE/GocacheUtilityOverTcp/client"
)

func main() {

	cli := c.NewClient(3000, "tcp", nil)

	cli.Connect()
	defer cli.Disconnect()

	cli.Set("BUSY", "intern")

	value, err := cli.Get("BUSY")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("value: %v\n", value)
}
```

Output:

```text
value: intern
```

---

## Client Constructor

Create a client using:

```go
cli := c.NewClient(3000, "tcp", nil)
```

Parameters:

| Parameter      | Type     | Description                              |
| -------------- | -------- | ---------------------------------------- |
| port           | int      | Server port                              |
| connectionType | string   | Network protocol (typically tcp)         |
| connObj        | net.Conn | Existing connection object (usually nil) |

---

## Example Operations

### SET

```go
cli.Set("username", "<User>")
```

### GET

```go
value, err := cli.Get("username")

if err == nil {
	fmt.Println(value)
}
```

Output:

```text
<User>
```

### Disconnect

```go
cli.Disconnect()
```

---

## Performance Benchmarks

### Sequential Performance (Single client multiple request) Before Auto Persistence and with SET response

| Test                                        | Result     |
| ------------------------------------------- | ---------- |
| 10,000 SET Operations                       | 3349.81 µs |
| 10,000 SET Operations (without server logs) | 2803.87 µs |
| Single SET Latency                          | 116.73 µs  |

---

### Throughput

| Clients | Operations per Client | Total Ops | Time   | Ops/sec |
| ------- | --------------------- | --------- | ------ | ------- |
| 100     | 1,000                 | 100,000   | 8.15s  | 12,267  |
| 10      | 10,000                | 100,000   | 16.08s | 6,216   |
| 1       | 100,000               | 100,000   | 31.36s | 3,186   |

---

### Persistence Impact (cleint * ops)

| Mode              | Time   | Ops/sec |
| ----------------- | ------ | ------- |
| DB Write Enabled  | 11.55s | 8,651   |
| DB Write Disabled (1 * 100000)| 13.57s | 7,365  |
| DB Write Disabled (10 * 10000)| 12.96s | 7710.55|
| DB Write Disabled (100 * 1000)| 7.07s  | 14,130 |

---

### Fire-and-Forget SET

| Clients | Total Ops | Time  | Ops/sec |
| ------- | --------- | ----- | ------- |
| 1       | 100,000   | 753ms | 132,750 |
| 10      | 100,000   | 221ms | 450,599 |
| 100     | 100,000   | 220ms | 454,347 |

---

### Fire-and-Forget GET

| Clients | Total Ops | Time  | Ops/sec |
| ------- | --------- | ----- | ------- |
| 100     | 100,000   | 737ms | 135,578 |
| 10      | 100,000   | 839ms | 119,103 |
| 1       | 100,000   | 2.77s | 35,991  |

---

### After Automatic Persistence

| Clients | Total Ops | Time  | Ops/sec |
| ------- | --------- | ----- | ------- |
| 100     | 100,000   | 219ms | 456,436 |
| 10      | 100,000   | 202ms | 494,695 |

> These benchmarks are intended for educational and comparative purposes and are not production-grade benchmark results.

---

## Typical Workflow

```text
1. Start Server
       ↓
2. Create Client
       ↓
3. Connect
       ↓
4. SET / GET Operations
       ↓
5. Automatic Persistence
       ↓
6. Automatic Compaction
       ↓
7. Disconnect
```

---

## Concepts Explored

This project helped deepen understanding of:

* TCP Networking
* Client-Server Architecture
* Goroutines
* Concurrent Programming
* RWMutex Synchronization
* JSON Serialization
* In-Memory Databases
* Persistence Mechanisms
* Storage Compaction
* Graceful Shutdown Handling
* Database Recovery
* Systems Programming in Go

---

## Comparison with Redis

| Feature            | Redis | GoCache Utility |
| ------------------ | ----- | --------------- |
| In-Memory Storage  | ✅     | ✅               |
| TCP Server         | ✅     | ✅               |
| Key-Value Store    | ✅     | ✅               |
| Persistence        | ✅     | ✅               |
| Multiple Clients   | ✅     | ✅               |
| Automatic Recovery | ✅     | ✅               |
| Pub/Sub            | ✅     | ❌               |
| Replication        | ✅     | ❌               |
| Clustering         | ✅     | ❌               |
| Transactions       | ✅     | ❌               |
| TTL Expiry         | ✅     | ❌               |

This project is intended as an educational implementation and is not intended to replace Redis.

---

## Future Improvements

* DELETE command
* Key expiration (TTL)
* Snapshot versioning
* Authentication
* Background replication
* REST API Gateway
* Binary protocol support
* Docker support
* Benchmark suite
* Unit tests
* Integration tests
* Cluster support

---

## Author

**Ankit Panchal**

Built to explore the foundations of Redis-style cache systems, persistence mechanisms, networking, concurrency, and systems programming in Go.
