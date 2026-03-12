# ---------- Build Stage ----------
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.19

WORKDIR /app

# Install certificates
RUN apk add --no-cache ca-certificates

# Copy compiled binary
COPY --from=builder /app/app .

# Run the application
CMD ["./app"]
