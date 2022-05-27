
build:
	docker-compose build

unittest:
	go test ./...

lint:
	go fmt ./...
	golangci-lint run --enable-all ./...

run:
	docker-compose up

test:
	docker-compose -f docker-compose.test.yml up --exit-code-from integration_tests

.PHONY: build unittest lint run test
