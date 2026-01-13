# URLchemy — A High-Performance URL Shortener (Go + Redis)

A high-performance, containerized URL shortener built with **Go 1.24** and **Redis**. This service utilizes cryptographically secure random generation and Base62 encoding to create short, shareable links with automatic expiration.



## Features
* **Base62 Encoding:** Generates clean, 6-character alphanumeric codes.
* **Secure Randomness:** Uses `crypto/rand` to ensure high entropy and low collision rates.
* **Auto-Expiration:** All links are automatically purged by Redis after a **24-hour TTL**.
* **Robust Orchestration:** Docker Compose configuration includes Redis health checks to ensure the application only starts when the database is ready.
* **Standard Library Driven:** Minimal external dependencies (uses `net/http` for routing).

---

## Tech Stack
* **Backend:** Go 1.24 (Alpine-based build)
* **Database:** Redis (Alpine-based image)
* **Environment Management:** `godotenv`
* **Containerization:** Docker & Docker Compose

---

## Configuration
The service is configured via environment variables defined in `docker-compose.yml`:

| Variable | Default | Description |
| :--- | :--- | :--- |
| `REDIS_ADDR` | `redis:6379` | Internal address of the Redis container |
| `APP_HOST_NAME` | `localhost:8000` | Hostname for generated short URLs |
| `APP_SCHEME` | `http` | Protocol (http/https) |

---

## API Endpoints

### 1. Create Short URL
* **Endpoint:** `POST /longurl`
* **Request Body:**
    ```json
    {
      "url": "[https://www.example.com/some/very/long/link](https://www.example.com/some/very/long/link)"
    }
    ```
* **Response (201 Created):**
    ```json
    {
      "url": "http://localhost:8000/aB3xYz"
    }
    ```

### 2. Redirect to Original
* **Endpoint:** `GET /{short_code}`
* **Response:** `302 Found` (Redirects to destination)

> [!IMPORTANT]
> **Validation Rule:** The service strictly requires the `url` to contain both a scheme (e.g., `https://`) and a hostname.

---

## Getting Started

### Prerequisites
* [Docker](https://www.docker.com/get-started)
* [Docker Compose](https://docs.docker.com/compose/install/)

### Deployment
1.  **Build and Start:**
    ```bash
    docker-compose up --build
    ```
    *The `url-shortner` service will wait for Redis to become "healthy" via its internal ping check before booting.*

2.  **Stop Service:**
    ```bash
    docker-compose down
    ```

---

## Project Structure

```text
.
├── main.go             # Application entry point & URL logic
├── Dockerfile          # Multi-stage Docker build (Golang 1.24-alpine)
├── docker-compose.yml  # Network and Service orchestration
├── go.mod              # Module definition
├── go.sum              # Dependency checksums
└── README.md           # Documentation
```

## Internal Logic

**Code Generation**: The service converts random bytes into a string using the alphabet: `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`.

**Collision Handling**: The server pings Redis to check if a code exists. If it does, it loops until a unique code is generated.

**Networking**: Containers communicate over a dedicated bridge network named app-network.