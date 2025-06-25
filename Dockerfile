# syntax=docker/dockerfile:1
FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache build-base

ENV CGO_ENABLED=1

COPY go.mod* go.sum* ./
RUN go mod download

COPY . ./
RUN go build -o gohttpx ./main.go

EXPOSE 8080

CMD ["./gohttpx"]
