default: pluvio-api

BUILD_TIME=$(shell date +%FT%T%z)
GIT_REVISION=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --symbolic-full-name --abbrev-ref HEAD)
GIT_DIRTY=$(shell git diff-index --quiet HEAD -- || echo "x-")

LDFLAGS=-ldflags "-s -X main.BuildStamp=$(BUILD_TIME) -X main.GitHash=$(GIT_DIRTY)$(GIT_REVISION) -X main.gitBranch=$(GIT_BRANCH)"

srcfiles = cmd/pluvio-api/main.go

testpackages = ./...

pluvio-api: $(srcfiles)
	go build -o bin/pluvio-api $(LDFLAGS) cmd/pluvio-api/main.go

docker:
	docker-compose up

docker-down:
	docker-compose down

docker-build:
	docker build $(DOCKER_BUILD_ARGS) -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

lint: 
	golangci-lint run --deadline 20m

test: 
	go test -count=1 $(testpackages)

integration-test: 
	@echo "Running integration tests, please ensure RMQ is running..."
	go test --tags=integration -count=1 $(testpackages)

test-verbose: 
	go test -v $(testpackages)

run-dev:
	CompileDaemon -command=bin/pluvio-api -build="make"

mod-tidy:
	go mod tidy

show-deps:
	go list -m all

clean:
	@rm -f bin/*

tidy: 
	go mod tidy
