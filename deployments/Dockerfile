# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /land-services ./cmd/land-services

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /land-services .
COPY configs/app.yaml ./configs/app.yaml

ENV LAND_CONFIG_FILE=/app/configs/app.yaml

EXPOSE 8199

ENTRYPOINT ["./land-services"]
