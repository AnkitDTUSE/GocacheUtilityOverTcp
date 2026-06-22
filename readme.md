# GoCache Utility Over TCP

A lightweight Redis-inspired key-value cache built in Go with TCP networking support.

GoCache Utility Over TCP is a simple client-server cache system that allows applications to store and retrieve key-value pairs over a TCP connection. The server maintains an in-memory cache for high-speed access while persisting data to disk using an append-only storage model.

This project was built to explore concepts behind distributed cache systems such as Redis, including networking, persistence, startup recovery, concurrent client handling, and log compaction.

---

## Features

* TCP-based client-server architecture
* In-memory key-value storage
* Persistent CSV-backed storage
* Append-only write strategy
* Automatic recovery on startup
* Manual log compaction
* JSON-based request protocol
* Concurrent client handling using Goroutines
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

The server keeps all active data in memory while persisting updates to disk.

Clients communicate with the server using JSON messages over TCP.

---

## How It Works

### In-Memory Storage

The cache uses:

```go
map[string]string
```

This provides near O(1) average lookup and insertion performance.

---

### TCP Communication

Clients send JSON requests to the server.

Example:

```json
{
  "cmd": "SET",
  "key": "name",
  "value": "YourName"
}
```

The server processes the request and sends a response back over the same TCP connection.

---

### Persistence

Every successful `SET` operation is appended to a CSV file.

Example:

```csv
path,C:
username,<User>
password,123456
```

This append-only strategy minimizes disk writes and mimics Redis-style AOF (Append Only File) persistence.

---

### Startup Recovery

When the server starts, it loads all records from `db.csv` and reconstructs the latest state of the cache.

Example:

```go
LoadData()
```

---

### Compaction

Since updates are appended rather than overwritten, duplicate key entries accumulate over time.

Before compaction:

```csv
user,John
user,Bob
user,Charlie
```

After compaction:

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
  "value": "YourName"
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
YourName
```

---

### COMPACT

Rewrite the database file using only the latest values stored in memory.

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

## Recommended Project Structure

When using this library, it is recommended to run the server as a separate process.

```text
my-project
├── server
│   ├── db.csv
│   ├── server.go
│   ├── go.mod
│   └── go.sum
│
├── client.go
├── go.mod
└── go.sum
```

The server should be started before any client attempts to connect.

---

## Starting the Server

Create `server/server.go`:

```go
package main

import (
	"fmt"

	s "github.com/AnkitDTUSE/GocacheUtilityOverTcp/server"
)

func main() {
	err := s.Start(3000, "tcp")

	if err != nil {
		fmt.Println("Error while starting server")
	}
}
```

Run the server:

```bash
cd server
go run .
```

The server will:

* Listen on port `3000`
* Create or load `db.csv`
* Recover previously stored data
* Accept multiple concurrent client connections

---

## Creating a Client

Create `client.go`:

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

	cli.Set("key", "asgdjkashdf5789")

	value, err := cli.Get("key")

	if err != nil {
		fmt.Println("error while fetching GET request")
		return
	}

	fmt.Println(value)

	cli.Compact()
}
```

Run:

```bash
go run client.go
```

---

## Client Constructor

The client package provides a constructor for creating client instances.

```go
cli := c.NewClient(3000, "tcp", nil)
```

Parameters:

| Parameter      | Type     | Description                                |
| -------------- | -------- | ------------------------------------------ |
| port           | int      | TCP port of the cache server               |
| connectionType | string   | Network protocol (typically `"tcp"`)       |
| connObj        | net.Conn | Existing connection object (usually `nil`) |

Example:

```go
cli := c.NewClient(3000, "tcp", nil)
```

This constructor simplifies initialization and follows idiomatic Go design patterns.

---

## Example Usage

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

	cli.Set("key", "asgdjkashdf5789")

	value, err := cli.Get("key")

	if err != nil {
		fmt.Println("error while fetching GET request")
		return
	}

	fmt.Println(value)

	cli.Compact()
}
```

Output:

```text
value: asgdjkashdf5789
Compaction Done
```

---

## Example Operations

### Writing Data

```go
cli.Set("username", "<User>")
cli.Set("password", "123456")
```

---

### Reading Data

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
5. COMPACT (Optional)
       ↓
6. Disconnect
```

---

## Concepts Explored

This project helped deepen understanding of:

* TCP Networking
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
* Systems Programming in Go

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
| Replication         | ❌     | ❌               |
| Clustering          | ❌     | ❌               |
| Transactions        | ❌     | ❌               |
| TTL Expiry          | ❌     | ❌               |

This project is intended as an educational implementation and is not meant to replace Redis.

---

## Author

**Ankit Panchal**

Built to explore the foundations of Redis-style cache servers, persistence mechanisms, networking, concurrency, and systems programming in Go.
