FROM golang:1.24.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o manpower ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root

COPY --from=builder /app/manpower .

ENTRYPOINT ["./manpower"]