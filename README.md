# Go-Redis: A High-Performance RESP-Compatible Key-Value Store

This is a minimalist, high-performance in-memory key-value store built from first principles using **Go**. It implements the **Redis Serialization Protocol (RESP)**, allowing it to communicate seamlessly with the official `redis-cli`.



## 🚀 Features

* **RESP Protocol Mastery:** Built a custom parser from scratch to handle Bulk Strings, Arrays, and Simple Strings.
* **Concurrency-Safe Storage:** Leverages Go’s `sync.RWMutex` to handle high-concurrency read/write operations without race conditions.
* **Persistence (AOF):** Implements **Append-Only File** logging to ensure data durability. Every write-operation is logged to disk for crash recovery.
* **Smart TTL (Time-To-Live):** Dynamic key expiration using Go's `time.AfterFunc` goroutines, with race-condition protection for updated keys.
* **Zero Dependencies:** Built entirely using Go’s standard library to understand the underlying mechanics of networking and memory.

## 🛠️ Technical Implementation

### Networking & Concurrency
The server uses a **Goroutine-per-connection** model, allowing it to scale to thousands of concurrent clients efficiently. Unlike the single-threaded event loop of the original Redis, this project leverages Go's runtime scheduler for parallelism.

### The Storage Engine
Data is stored in a hash map guarded by a **Read-Write Mutex**. This allows multiple concurrent readers but ensures exclusive access for writers, optimizing for the typical read-heavy workload of a cache.

### Persistence Mechanism
Inspired by Redis's AOF, the system records every state-changing command. On startup, the server "replays" the log file to reconstruct the in-memory state.



## 🚦 Getting Started

1. **Clone the repository:**
   ```bash
   git clone [https://github.com/yourusername/go-redis.git](https://github.com/yourusername/go-redis.git)
   cd go-redis