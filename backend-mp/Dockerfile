# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/seed ./cmd/seed

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/bin/migrate .
COPY --from=builder /app/bin/seed .

EXPOSE 8080

CMD ["./server"]
