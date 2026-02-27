# San API - Simple Go Backend Boilerplate

This project is a simple, educational boilerplate for building RESTful APIs using **Go (Golang)**. It demonstrates clean architecture, database management, and asynchronous processing patterns.

## 📚 Educational Goals

By exploring this codebase, you will learn:
- **Clean Architecture**: How to separate concerns into Handlers, Services, and Repositories.
- **Dependency Injection**: How to wire components together for better testability and modularity.
- **Database Management**: Using **SQLC** for type-safe SQL queries and **Golang-Migrate** for schema versioning.
- **Asynchronous Processing**: Using **RabbitMQ** for reliable background task processing (e.g., email sending).
- **Configuration**: Managing environment variables with **Viper**.
- **Documentation**: Auto-generating API docs with **Swagger (Swag)**.
- **Containerization**: running the application, worker, database, and message broker with **Docker Compose**.

## 🚀 Tech Stack

- **Language**: Go 1.24+
- **Framework**: [Gin](https://github.com/gin-gonic/gin) (HTTP Web Framework)
- **Database**: PostgreSQL
- **ORM/Query Builder**: [SQLC](https://sqlc.dev/) (Generate type-safe Go from SQL)
- **Message Broker**: [RabbitMQ](https://www.rabbitmq.com/) (Background jobs)
- **Migration**: [Golang-Migrate](https://github.com/golang-migrate/migrate)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Documentation**: [Swag](https://github.com/swaggo/swag)
- **Hot Reload**: [Air](https://github.com/air-verse/air)

## ✨ Key Features & Highlights

### 🔐 Authentication & Security
- **JWT Authentication**: robust access and refresh token implementation.
- **Owner-Only Deletion**: Strict authorization policy ensuring users can only delete their own posts.
- **Protected Routes**: Middleware-guarded endpoints for sensitive operations.

### � Background Jobs & Email
- **Worker Service**: Dedicated worker service running separately from the API server for better scalability.
- **RabbitMQ Integration**: Reliable task queuing and distribution.
- **Email Verification**: Asynchronous email sending flow for new user registration.
- **Mailtrap Support**: Easy development testing with Mailtrap integration (or custom SMTP).

### �� Post Management
- **Full CRUD**: comprehensive operations for blog posts.
- **Abstract Field**: A dedicated summary field separate from the main content. Perfect for:
    - SEO meta descriptions
    - Blog list previews (cards)
    - Social media sharing snippets
- **Rich Metadata**: Support for tags, publication status, and location data.

### 📊 Observability
- **Custom Logger Middleware**: Tracks and logs:
    - HTTP Status Codes
    - Response Time (Latency)
    - Client IP Addresses
    - HTTP Methods & Paths

## 📂 Project Structure

```
.
├── cmd/
│   ├── server/       # API Server entry point (main.go)
│   └── worker/       # Background Worker entry point (main.go)
├── internal/
│   ├── config/       # Configuration loading logic
│   ├── db/           # Database connection and migration logic
│   │   ├── migration/# SQL migration files
│   │   ├── query/    # SQL queries for SQLC
│   │   └── sqlc/     # Generated Go code from SQLC
│   ├── handler/      # HTTP Handlers (Controllers) - Handle requests/responses
│   ├── service/      # Business Logic Layer - Core application logic
│   ├── server/       # Server setup and routing
│   └── worker/       # Worker logic (Distributor & Processor)
├── pkg/              # Public libraries (Logger, Utils, Mail)
├── api/              # API definitions (DTOs)
├── docs/             # Swagger documentation files
├── env/              # Environment configuration files
└── Makefile          # Development commands
```

## 🛠 Prerequisites

- **Go** (version 1.22 or higher)
- **Docker** & **Docker Compose**
- **Make** (optional, for running Makefile commands)

## ⚡️ Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/trungnhann/san.git
cd san
```

### 2. Configure Environment

Copy the development environment file to `.env`:

```bash
cp env/development.env .env
```

If you want to test email sending, update the following variables in `.env` with your **Mailtrap** credentials:

```env
EMAIL_SENDER_HOST=sandbox.smtp.mailtrap.io
EMAIL_SENDER_PORT=2525
EMAIL_SENDER_USERNAME=<your-mailtrap-username>
EMAIL_SENDER_PASSWORD=<your-mailtrap-password>
EMAIL_SENDER_ADDRESS=no-reply@san.dev
```

### 3. Run with Docker (Recommended)

This will start the API server, Worker, PostgreSQL, Redis, RabbitMQ, and MinIO.

```bash
docker-compose up --build -d
```

- **API Server**: `http://localhost:3001`
- **RabbitMQ UI**: `http://localhost:15672` (User/Pass: `guest`/`guest`)
- **MinIO Console**: `http://localhost:9001` (User/Pass: `minioadmin`/`minioadmin`)

### 4. Run Locally (Alternative)

Ensure you have PostgreSQL, Redis, and RabbitMQ running locally.

```bash
# Install dependencies
go mod download

# Run the API server
go run cmd/server/main.go

# Run the Worker (in a separate terminal)
go run cmd/worker/main.go
```

## 📖 API Documentation

Once the server is running, you can access the interactive Swagger documentation at:

👉 **[http://localhost:3001/swagger/index.html](http://localhost:3001/swagger/index.html)**

## 🔧 Development Commands

The project includes a `Makefile` to simplify common tasks:

| Command | Description |
|---------|-------------|
| `make dev` | Run the server with hot-reload (requires Air) |
| `make migration name=create_users` | Create a new database migration file |
| `make sqlc` | Generate Go code from SQL queries (requires SQLC) |
| `make swagger` | Generate Swagger documentation from code comments |

### Example Workflow

1.  **Add a new feature**:
    - Create a migration: `make migration name=add_posts_table`
    - Write SQL in `internal/db/migration/`
    - Write query in `internal/db/query/`
    - Generate code: `make sqlc`
    - Implement Service and Handler logic.

2.  **Update API Docs**:
    - Add comments to your handler functions.
    - Run `make swagger`.

## 💡 Key Concepts Explained

### 1. Interface Segregation
We define interfaces where they are **used** (Consumer), not where they are implemented.
- `internal/service/interfaces.go`: Defines repository interfaces needed by the Service.
- `internal/handler/interfaces.go`: Defines use-case interfaces needed by the Handler.

### 2. SQLC vs ORM
Instead of a heavy ORM like GORM, we use **SQLC**. You write raw SQL, and it generates type-safe Go code. This gives you the performance of raw SQL with the safety of a compiled language.

### 3. Background Processing (RabbitMQ)
Long-running tasks (like sending emails) are offloaded to **RabbitMQ**.
- **Distributor**: The API server pushes tasks to the queue.
- **Processor**: The Worker service consumes tasks from the queue and executes them reliably.
