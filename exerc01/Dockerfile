FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o configparser ./main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/configparser .

ENTRYPOINT ["/app/configparser"]