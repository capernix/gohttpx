# GoHTTPx - Modern Go REST API with CLI

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)
[![SQLite](https://img.shields.io/badge/Database-SQLite-green.svg)](https://sqlite.org)
[![Nginx](https://img.shields.io/badge/Proxy-Nginx-green.svg)](https://nginx.org)

A full-stack Go web application demonstrating modern backend development practices with REST API, database persistence, containerization, and a command-line interface.

## Features

### REST API Server
- **Modern Go HTTP Server** - Built with Go 1.24's enhanced HTTP routing
- **RESTful Architecture** - Clean API design following REST principles
- **SQLite Database** - Persistent data storage with automatic schema migration
- **Graceful Shutdown** - Proper signal handling and cleanup
- **CORS Support** - Cross-origin resource sharing enabled
- **JSON Responses** - Structured error handling and response formatting

### Infrastructure & DevOps
- **Docker Containerization** - Multi-stage Docker builds for production
- **Docker Compose** - Complete development environment setup
- **Nginx Reverse Proxy** - Load balancing and static file serving
- **Persistent Volumes** - Database data persistence across container restarts
- **Health Checks** - Container health monitoring

### Command Line Interface
- **Cobra CLI Framework** - Professional command-line interface
- **Multiple Commands** - User management, note operations, and configuration
- **Configuration Management** - Flexible endpoint and authentication setup

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Client      │    │      Nginx      │    │    Go Server    │
│   (Browser/CLI) │◄──►│ Reverse Proxy   │◄──►│   (Port 8080)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                │                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Static Files   │    │  SQLite Database│
                       │   (HTML/CSS)    │    │    (Persistent) │
                       └─────────────────┘    └─────────────────┘
```

## Technology Stack

### Backend
- **Language**: Go 1.24
- **HTTP Framework**: Native Go HTTP server with enhanced routing
- **Database**: SQLite 3 with CGO support
- **CLI Framework**: Cobra (spf13/cobra)

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Nginx
- **Database Driver**: mattn/go-sqlite3

### Development Tools
- **Package Management**: Go Modules
- **Container Orchestration**: Docker Compose
- **Version Control**: Git

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.24+ (for local development)
- Git

### 1. Clone and Run with Docker
```bash
git clone https://github.com/capernix/gohttpx.git
cd gohttpx
docker-compose up --build -d
```

### 2. Access the Application
- **Web Interface**: http://localhost
- **API Endpoints**: http://localhost/api/*
- **Direct API**: http://localhost:8080/* (development)

### 3. Test the API
```bash
# Create a user
curl -X POST http://localhost/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Get all users
curl http://localhost/users

# Create a note
curl -X POST http://localhost/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"Meeting Notes","content":"Discuss project roadmap","user_id":1}'

# Get all notes
curl http://localhost/notes
```

## API Documentation

### Users Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/users` | Create a new user |
| GET    | `/users` | List all users |
| GET    | `/users/{id}` | Get user by ID |
| DELETE | `/users/{id}` | Delete user by ID |

### Notes Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/notes` | Create a new note |
| GET    | `/notes` | List all notes |
| GET    | `/notes/{id}` | Get note by ID |
| DELETE | `/notes/{id}` | Delete note by ID |

### Request/Response Examples

#### Create User
```json
// POST /users
{
  "name": "John Doe",
  "email": "john@example.com"
}

// Response
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

#### Create Note
```json
// POST /notes
{
  "title": "Project Planning",
  "content": "Plan the next sprint activities",
  "user_id": 1
}

// Response
{
  "status": "success",
  "data": {
    "id": 1,
    "title": "Project Planning",
    "content": "Plan the next sprint activities",
    "user_id": 1,
    "created_at": "2025-01-15T10:35:00Z"
  }
}
```

## CLI Usage

### Build and Install CLI
```bash
cd cli
go build -o gohttpx-cli .
```

### CLI Commands
```bash
# Configure API endpoint
./gohttpx-cli config set-url http://localhost:8080

# User management
./gohttpx-cli user create --name "Jane Doe" --email "jane@example.com"
./gohttpx-cli user list
./gohttpx-cli user show --id 1

# Note management
./gohttpx-cli note create --title "Meeting" --content "Team standup" --user-id 1
./gohttpx-cli note list
./gohttpx-cli note show --id 1
```

## Development

### Local Development Setup
```bash
# Clone repository
git clone https://github.com/capernix/gohttpx.git
cd gohttpx

# Install dependencies
go mod download

# Run database migrations
go run main.go

# Start development server
go run main.go
```

### Environment Variables
```bash
# Database configuration
DB_PATH=./data/gohttpx.db

# Server configuration
PORT=8080
HOST=localhost
```

### Building from Source
```bash
# Build the server
go build -o gohttpx main.go

# Build the CLI
cd cli && go build -o gohttpx-cli .
```

## Docker Deployment

### Production Build
```bash
# Build and start services
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Docker Services
- **gohttpx**: Go API server (port 8080)
- **nginx**: Reverse proxy and static file server (port 80)

### Data Persistence
- SQLite database stored in `./data/` directory
- Persistent across container restarts
- Automatic backup recommended for production

## Testing

### Manual API Testing
```bash
# Health check
curl http://localhost:8080/users

# Create test data
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com"}'
```

### Load Testing
```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost/users

# Using curl for concurrent requests
for i in {1..10}; do
  curl -X POST http://localhost/users \
    -H "Content-Type: application/json" \
    -d '{"name":"User'$i'","email":"user'$i'@example.com"}' &
done
```

## Security Features

- **Input Validation**: JSON request validation
- **SQL Injection Prevention**: Parameterized queries
- **CORS Configuration**: Controlled cross-origin access
- **Graceful Error Handling**: No sensitive data exposure
- **Container Security**: Non-root user in Docker containers

## Performance Considerations

- **Connection Pooling**: SQLite connection management
- **Graceful Shutdown**: Proper cleanup on termination
- **Static File Caching**: Nginx handles static assets
- **Gzip Compression**: Reduced bandwidth usage
- **Keep-Alive Connections**: HTTP connection reuse

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Your Name**
- GitHub: [@capernix](https://github.com/capernix)
- Project: [GoHTTPx](https://github.com/capernix/gohttpx)

---

*Built with ❤️ using Go, Docker, and modern development practices*
