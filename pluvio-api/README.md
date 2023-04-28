# Pluvio API

The backend API for pluvio functionality.

## 1. Setup

### Dependencies

- Docker
- Go (you need this to build the binaries on your local machine)

### 1.1 For Docker Setup:

This is the suggested way for development as it will autobuild/reload when changing
code.

```
docker-compose build
docker-compose up
```

You can now execute requests at port `7080` at `/api/v1/seating_solution` using `POST`.

For example:

```
curl -d '{"rows": [4,4], "groups": [2,3]}' -H "Content-Type: application/json" -H "Authorization: Token change-in-prod" -X POST http://localhost:7080/api/v1/seating_solution
```

### 1.2 For local setup:

```
go build -o ./bin/pluvio-api cmd/pluvio-api/main.go
./bin/pluvio-api
```

### 1.3 Testing:

```
go test ./...
```

### 1.4 Linting:

We use golangci-lint

```
golangci-lint run
```