# Orbit WebSocket Chat Service

Orbit is a high-performance WebSocket messaging service built with Go, following the **Hexagonal Architecture** (Ports and Adapters) pattern. It serves as a real-time communication layer, managing active user sessions, message persistence, and asynchronous broadcasting.

## 🚀 Features

- **Real-time Messaging**: Low-latency communication using WebSockets.
- **Asynchronous Broadcasting**: Utilizes **RabbitMQ** to distribute messages across multiple server nodes.
- **User-to-User Routing**: Intelligent message routing based on Recipient IDs.
- **Session Management**: Tracks active connections and notifies external systems about user status.
- **Persistence**: Automatically saves all incoming messages to a Core API.
- **Security**: JWT Authentication (RS256) for all connections.
- **Observability**: Structured logging with Trace IDs for request tracking.

## 🏗️ Architecture

The project follows a modular hexagonal design:

- **Core/Entities**: Clean business models (`SessionMessage`, `MessageReceived`).
- **Core/Ports**: Interface definitions for services (User Notifications, Message Store).
- **Core/Services**: Business logic implementations with Factory support for **Mock/HTTP** switching.
- **Adapters/Inbound**:
  - **WebSocket Handlers**: Manages the connection lifecycle.
  - **Middlewares**: Auth (JWT) and Logger.

## 🛠️ Tech Stack

- **Go**: Primary language.
- **RabbitMQ**: Message bus for scaling and delivery.
- **Coder/Websocket**: Robust WebSocket implementation.
- **Docker Compose**: For local infrastructure.

## 🚦 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### 1. Infrastructure Setup
Spin up RabbitMQ using the provided configuration:
```bash
docker-compose up -d
```

### 2. Configuration
Create a `.env.dev` file inside `.envs/` (or use existing ones) with the following variables:
```bash
PORT=8003
CORE_API_URL=mock # Change to real URL for production
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=orbit.chat
JWT_PUBLIC_KEY=-----BEGIN PUBLIC KEY----- ...
```

### 3. Run the Server
```bash
go run cmd/main.go
```

### 4. Testing with CLI Client
We provide a colored CLI tool to simulate multiple users:
```bash
# Compile the tool
go build -o ws-test cmd/test-client/main.go

# Run for different users (requires test_users.json)
./ws-test user_1
./ws-test user_2
```

## 🧪 Development & Mocks

Orbit is designed to be independent of external APIs during development. By setting `CORE_API_URL=mock`, the system will use internal mocks for:
- User Notifications (Log only).
- Message Storage (Log only).

This allows frontend or core development without needing all services running.

---
*Built with ❤️ by the Jalo Team.*
