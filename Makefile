
build:
	docker-compose build

unittest:
	go test ./...

lint:
	go fmt ./...
	golangci-lint run --enable-all ./...

run:
	docker-compose up

.PHONY: build unittest lint run
