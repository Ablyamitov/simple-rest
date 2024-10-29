
FROM golang:1.23 AS builder
LABEL authors="ToTheMoon"

WORKDIR /app

COPY . .
#toolchain
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simple-rest ./cmd/app


FROM alpine:latest AS migrate-installer

#migrate
RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz -o /tmp/migrate.tar.gz \
    && tar -xzf /tmp/migrate.tar.gz -C /tmp/ \
    && mv /tmp/migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate \
    && rm -rf /tmp/migrate.tar.gz /tmp/migrate


FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/simple-rest .
COPY --from=migrate-installer /usr/local/bin/migrate /usr/local/bin/migrate

COPY --from=builder /app/configs ./configs
COPY --from=builder /app/internal/store/db/migrations ./internal/store/db/migrations
ENV CONFIG_PATH=/app/configs/config.yaml

EXPOSE 8080

CMD ["./simple-rest"]
