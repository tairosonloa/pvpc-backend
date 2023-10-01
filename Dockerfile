FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 go build -o /app/bin/http ./cmd/http/
RUN CGO_ENABLED=0 go build -o /app/bin/migrate ./cmd/migrate/

FROM alpine:latest

RUN apk add --no-cache tzdata

EXPOSE 8080

COPY scripts/start.sh /app/start.sh
COPY internal/platform/storage/postgresql/migrations /app/migrations
COPY --from=builder /app/bin/http /app/bin/http
COPY --from=builder /app/bin/migrate /app/bin/migrate

CMD ["/app/start.sh"]