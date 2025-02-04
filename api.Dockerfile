# Build container
FROM golang:alpine AS builder

RUN apk update && apk upgrade && apk add git zlib-dev gcc musl-dev

COPY . /go/src/github.com/TicketsBot/export
WORKDIR /go/src/github.com/TicketsBot/export

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -o api cmd/api/main.go

# Prod container
FROM alpine:latest

RUN apk update && apk upgrade && apk add curl

COPY --from=builder /go/src/github.com/TicketsBot/export/api /srv/api/api

RUN chmod +x /srv/api/api

RUN adduser container --disabled-password --no-create-home
USER container
WORKDIR /srv/api

CMD ["/srv/api/api"]