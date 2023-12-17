# syntax=docker/dockerfile:1.4

# Start from the official Golang image to build the binary file
FROM golang:1.20 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd /app/cmd

COPY ./pkg /app/pkg

WORKDIR /app/cmd/repositories-service

RUN CGO_ENABLED=0 GOOS=linux go build -o ./main .

# Multi stage build
FROM alpine:latest as runner

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/cmd/repositories-service/main .

COPY --link setup.sh /root/

EXPOSE 8080

CMD ["./main"]
