SHELL=/bin/bash

.PHONY: openapi-build
openapi-build:
	go generate cmd/run-service/main.go

.PHONY: go-tidy
go-tidy:
	go mod tidy 

.PHONY: run-docker
run-docker:
	docker-compose up -d --build

.PHONY: stop-docker
stop-docker:
	docker-compose down -v  
	docker system prune -a -f --volumes

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: run 
run:
	go run cmd/run-service/main.go

.PHONY: gofmt
gofmt:
	go fmt  ./...

.PHONY: test
test:
	go test -v ./...
