# GoCache Utility Over TCP

A lightweight Redis-inspired key-value cache built in Go with TCP networking support.

This project is the successor to my earlier local cache implementation and extends it into a networked client-server architecture. The cache supports multiple TCP clients, persistent storage, append-only logging, startup recovery, and manual log compaction.

The goal of this project is to understand how distributed cache systems and in-memory databases communicate over the network while maintaining persistence and high-speed data access.

---

## Features

* TCP-based client-server architecture
* In-memory key-value storage
* Persistent CSV-backed storage
* Append-only write strategy
* Automatic recovery on server startup
* Manual log compaction
* JSON-based request protocol
* Concurrent client handling using Goroutines
* Fast lookups using Go maps

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

The server maintains the cache in memory while persisting writes to disk.

Clients communicate with the server using JSON messages over TCP.

---

## How It Works

### In-Memory Storage

The server stores data in a Go map:

```go
map[string]string
```

This allows near O(1) average lookup and insertion performance.

---

### TCP Communication

Clients send JSON requests to the server:

```json
{
  "cmd": "SET",
  "key": "name",
  "value": "Ankit"
}
```

The server processes the request and sends a response back over the same TCP connection.

---

### Persistence

Every successful `SET` operation is appended to a CSV file.

Example:

```csv
Vineet,panchal
Vineet,singh
Vineet,Kumar
Sanjay,panchal
```

This append-only strategy minimizes disk writes and mimics concepts used in Redis AOF (Append Only File) persistence.

---

### Startup Recovery

When the server starts, it loads data from `db.csv`.

Example:

```go
LoadData()
```

The latest value for every key is reconstructed by replaying the log file.

---

### Compaction

Since updates are appended rather than overwritten, duplicate key entries accumulate over time.

Before compaction:

```csv
user,John
user,Bob
user,Charlie
```

After running:

```text
COMPACT
```

The database file becomes:

```csv
user,Charlie
```

This reduces storage size and improves recovery time.

---

## Supported Commands

### SET

Store a key-value pair.

Request:

```json
{
  "cmd": "SET",
  "key": "name",
  "value": "Ankit"
}
```

Response:

```text
OK
```

---

### GET

Retrieve the latest value associated with a key.

Request:

```json
{
  "cmd": "GET",
  "key": "name"
}
```

Response:

```text
Ankit
```

---

### COMPACT

Rewrite the database file using only the latest values currently stored in memory.

Request:

```json
{
  "cmd": "COMPACT"
}
```

Response:

```text
Compaction Complete
```

---

## Example Usage

### Starting the Server

```go
go s.Start(3000, "tcp")
```

---

### Creating a Client

```go
cli := c.Client{
    Port:           3000,
    ConnectionType: "tcp",
}
```

---

### Connecting

```go
cli.Connect()
```

---

### Writing Data

```go
cli.Set("Vineet", "panchal")
cli.Set("Vineet", "singh")
cli.Set("Vineet", "Kumar")
```

---

### Reading Data

```go
value, _ := cli.Get("Vineet")
fmt.Println(value)
```

Output:

```text
Kumar
```

---

### Compacting Storage

```go
cli.Compact()
```

---

### Disconnecting

```go
cli.Disconnect()
```

---

## Project Structure

```text
.
├── client
│   └── client.go
│
├── server
│   ├── server.go
│   ├── cache.go
│   └── persistence.go
│
├── db.csv
├── main.go
└── README.md
```

---

## Concepts Explored

This project helped deepen understanding of:

* TCP Networking in Go
* Client-Server Architecture
* JSON Serialization
* Goroutines
* Concurrent Connection Handling
* In-Memory Databases
* Persistence Mechanisms
* Append-Only Logging
* Log Compaction
* File Handling
* Database Recovery
* System Design Fundamentals

---

## Comparison with Redis

| Feature             | Redis | GoCache Utility |
| ------------------- | ----- | --------------- |
| In-Memory Storage   | ✅     | ✅               |
| TCP Server          | ✅     | ✅               |
| Key-Value Store     | ✅     | ✅               |
| Persistence         | ✅     | ✅               |
| Append-Only Logging | ✅     | ✅               |
| Log Compaction      | ✅     | ✅               |
| Multiple Clients    | ✅     | ✅               |
| Pub/Sub             | ✅     | ❌               |
| Replication         | ✅     | ❌               |
| Clustering          | ✅     | ❌               |
| Transactions        | ✅     | ❌               |
| TTL Expiry          | ✅     | ❌               |

This project is not intended to replace Redis. It is an educational implementation focused on understanding the internals of networked cache systems.

---

## Future Improvements

* DELETE command
* Key expiration (TTL)
* Snapshot persistence
* Background compaction
* Authentication
* Benchmark suite
* Binary protocol support
* Replication between servers
* REST API Gateway
* Docker support
* Unit and integration tests

---

## Installation

Clone the repository:

```bash
git clone https://github.com/AnkitDTUSE/GocacheUtilityOverTcp.git
```

```bash
cd GocacheUtilityOverTcp
```

Install dependencies:

```bash
go mod tidy
```

Run:

```bash
go run .
```

---

## Author

**Ankit Panchal**

Built as a systems programming project to explore the foundations of Redis-style cache servers, persistence mechanisms, and networked database design in Go.
