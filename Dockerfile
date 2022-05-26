# Environment
FROM golang:1.18.2 as build-env

RUN mkdir -p /opt/go.audit
WORKDIR /opt/go.audit
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/service/audit_api

# Release
FROM alpine:latest
COPY --from=build-env /opt/service/audit_api /bin/audit_api
RUN mkdir /etc/audit
COPY example_config.yaml /etc/audit/config.yaml
ENTRYPOINT ["/bin/audit_api", "--config", "/etc/audit/config.yaml"]
