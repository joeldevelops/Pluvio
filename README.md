# Pluvio - Rain Data Reporting for the Global South

![Pluvio Logo](Website/logo-transparent.png)

- [Pluvio - Rain Data Reporting for the Global South](#pluvio---rain-data-reporting-for-the-global-south)
	- [Folder Structure](#folder-structure)
	- [Pluvio API](#pluvio-api)
	- [1. Setup](#1-setup)
		- [Dependencies](#dependencies)
		- [1.1 For Docker Setup (NOT IMPLEMENTED YET):](#11-for-docker-setup-not-implemented-yet)
		- [1.2 For local setup:](#12-for-local-setup)
		- [1.3 Testing:](#13-testing)
		- [1.4 Linting:](#14-linting)
	- [2. Rain API](#2-rain-api)
		- [2.1 POST /api/v1/rain](#21-post-apiv1rain)
			- [Request](#request)
			- [200 Response](#200-response)
			- [400 Response](#400-response)
			- [403 Response](#403-response)
			- [429 Response](#429-response)
		- [2.2 GET /api/v1/rain/:timeRange](#22-get-apiv1raintimerange)
			- [Request](#request-1)
			- [Response](#response)
	- [3. User API](#3-user-api)
		- [3.1 POST /api/v1/user](#31-post-apiv1user)
			- [Request](#request-2)
			- [200 Response](#200-response-1)
			- [409 Response](#409-response)
			- [500 Response](#500-response)


## Folder Structure

```
.
├── Makefile
├── Procfile
├── README.md
├── VXML - Contains the VXML files which are copied into voxeo
│   ├── PourleMali.wav
│   ├── pluvio.xml
│   └── pluvio_prototype.xml
├── api - Go files
│   ├── api.go
│   └── handler.go
├── cmd - Contains the main.go files for the binaries
│   └── pluvio-api
│       └── main.go
├── go.mod
├── go.sum
└── mdb - Contains the MongoDB files for the application
    ├── models.go
    └── mongodb.go
```

## Pluvio API

The backend API for pluvio functionality. The Website/Frontend code is available here: [lds19/Pluvio](https://github.com/lds19/Pluvio)

## 1. Setup

### Dependencies

- Docker
- Go (you need this to build the binaries on your local machine)

### 1.1 For Docker Setup (NOT IMPLEMENTED YET):

TODO

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

## 2. Rain API

When making requests to the API, keep in mind that it can return either XML or JSON depending on the `Accept` header.
```
curl -H "Accept: application/json" -X GET http://localhost:4242/api/v1/rain/day?location=Mali
curl -H "Accept: application/xml" -X GET http://localhost:4242/api/v1/rain/day?location=Mali
```

### 2.1 POST /api/v1/rain

#### Request

```
curl -d '{"location": "Mali", "amount": 10, "phoneNumber":"+31612345678"}' -H "Content-Type: application/json" -X POST http://localhost:4242/api/v1/rain
```

#### 200 Response

```xml
<?xml version="1.0" ?>
<vxml version="2.1">
	<form>
		<block>
			<prompt>Thank you for your report!</prompt>
		</block>
	</form>
</vxml>
```

#### 400 Response
Response if no `phoneNumber` is provided and `USE_PHONE_AUTH` is set to `true` in the `.env` file.

```xml
<?xml version="1.0" ?>
<vxml version="2.1">
	<form>
		<block>
			<prompt>No phone number provided, please call back and try again.</prompt>
		</block>
	</form>
</vxml>
```

#### 403 Response
Response if  `USE_PHONE_AUTH` is set to `true` in the `.env` file and the user is not in the DB.

```xml
<?xml version="1.0" ?>
<vxml version="2.1">
	<form>
		<block>
			<prompt>You are not authorized to use this service.</prompt>
		</block>
	</form>
</vxml>
```

#### 429 Response
Reponse if `USE_PHONE_AUTH` is set to `true` in the `.env` file and you have already made a report for the last 24 hours.

```xml
<?xml version="1.0" ?>
<vxml version="2.1">
	<form>
		<block>
			<prompt>Sorry, you have reached the maximum number of reports for today.</prompt>
		</block>
	</form>
</vxml>
```

### 2.2 GET /api/v1/rain/:timeRange

#### Request

```
curl -X GET "http://localhost:4242/api/v1/rain/:timeRange?location=<location>"
```
Where:
- `:timeRange` is one of `day`, `week`, `month`
- (Optional) `location` is one of `Mali`, `(Burkina) Faso`. Leaving this out will return data from all locations.

#### Response

```xml
<?xml version="1.0" ?>
<vxml version="2.1">
	<form>
		<block>
			<prompt>In the past [day|week|month] it rained 100 milliliters</prompt>
		</block>
	</form>
</vxml>
```


## 3. User API

### 3.1 POST /api/v1/user

#### Request

```
curl -d '{"name": "Joel", "phone": "+31612345678"}' -H "Content-Type: application/json" -X POST http://localhost:4242/api/v1/user
```


#### 200 Response

```
Created user with ID: 645bb964c8dfb30099937009
```

#### 409 Response

```
User already exists
```

#### 500 Response

```
Internal server error
```