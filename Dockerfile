FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .
COPY go.mod .
RUN go mod download

RUN go build -o peroblogo .

# create a new slim container
FROM debian:stretch-slim

COPY ./wait-for.sh /
COPY ./conf /conf

COPY --from=builder /build/peroblogo /

RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		netcat; \
        chmod 755 wait-for.sh

# ENTRYPOINT ["./peroblogo", "-env=prod"]
