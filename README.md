# SU Server

Backend server for Mae Fah Luang University Student Union services, built with Go and PostgreSQL.

> 🚧 This project is currently under active development.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **Architecture:** Handler → Service → Repository pattern

## Project Structure

```
su-server/
├── cmd/
│   └── main.go               # Entry point
├── config/
│   └── database.go           # Database connection & config
├── db/
│   └── migrations/
│       └── 001_init.sql      # Initial schema
├── internal/
│   ├── handler/              # HTTP handlers (request/response)
│   ├── middleware/           # HTTP middleware
│   ├── model/                # Data models
│   ├── repository/           # Database queries
│   └── service/              # Business logic
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Features

### Events
- `GET /events` — List all events
- `GET /events/:id` — Get event by ID
- `POST /events` — Create a new event
- `PUT /events/:id` — Update an event
- `DELETE /events/:id` — Delete an event

> More services coming soon.

## Getting Started

### Prerequisites

- [Go 1.26+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) & Docker Compose

### Run with Docker

```bash
docker-compose up --build
```

### Run locally

```bash
# Start PostgreSQL
docker-compose up -d db

# Run the server
make run
```

### Database Migrations

```bash
make migrate
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `secret` |
| `DB_NAME` | Database name | `su_server` |
| `PORT` | Server port | `8080` |

## License

See [LICENSE](./LICENSE) for details.
