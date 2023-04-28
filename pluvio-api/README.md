# Pluvio API

The backend API for pluvio functionality.

## 1. Setup

### Dependencies

- Docker
- Go (you need this to build the binaries on your local machine)

### 1.1 For Docker Setup (NOT IMPLEMENTED YET):

This is the suggested way for development as it will autobuild/reload when changing
code.

```
docker-compose build
docker-compose up
```

You can now execute requests at port `4242` at `/api/v1/rain` using `POST` or `GET`.

For example:

```
curl -d '{"location": "Mali", "amount": 10}' -H "Content-Type: application/json" -X POST http://localhost:4242/api/v1/rain
```

### 1.2 For local setup:

```
cat .env.example > .env
go mod download
make
./bin/pluvio-api
```

Be sure to reach out to Joel to get the connection uri for the database. The other defaults should be fine.

### 1.3 Testing:

```
go test ./...
```

### 1.4 Linting:

We use golangci-lint

```
golangci-lint run
```

## 2. API

### 2.1 POST /api/v1/rain

#### Request

```
curl -d '{"location": "Mali", "amount": 10}' -H "Content-Type: application/json" -X POST http://localhost:4242/api/v1/rain
```

#### Response

```
"5f9b1b7b9d9b7b0001b9d9b7"
```

### 2.2 GET /api/v1/rain/day

#### Request

```
curl -X GET http://localhost:4242/api/v1/rain/day
```

#### Response

```
"100mm rain for the past day"
```

### 2.3 GET /api/v1/rain/week

#### Request

```
curl -X GET http://localhost:4242/api/v1/rain/week
```

#### Response

```
"100mm rain for the past week"
```

### 2.4 GET /api/v1/rain/month

#### Request

```
curl -X GET http://localhost:4242/api/v1/rain/month
```

#### Response

```
"100mm rain for the past month"
```