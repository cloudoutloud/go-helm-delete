# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o helm-delete main.go

# Stage 2: Create a lightweight image
FROM alpine:latest

# Create minimal, non-root user
RUN addgroup -S helmgroup && adduser -S helmuser -G helmgroup

RUN apk add --no-cache helm

WORKDIR /app

COPY --from=builder /app/helm-delete .

RUN chmod +x /app/helm-delete
RUN chown -R helmuser:helmgroup /app

USER helmuser

ENTRYPOINT ["./helm-delete"]