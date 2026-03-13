## Build Stage
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/main.go


## Runtime Stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    e2fsprogs \
    xfsprogs \
    util-linux \
    blkid

COPY --from=builder /app/app .

ENTRYPOINT ["./app"]