# Environment
FROM golang:1.18.2

RUN mkdir -p /opt/goaudit_tests
WORKDIR /opt/goaudit_tests
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD go test -v ./...