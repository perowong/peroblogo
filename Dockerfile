FROM golang:1.19 AS builder

WORKDIR /build

ADD .netrc /root/

ADD .keys/ /root/.ssh/
RUN chmod 700 -R /root/.ssh/
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/agent-chatee.insteadOf https://github.com/agent-chatee

ENV GOPROXY https://goproxy.cn
ENV GOPRIVATE github.com/agent-chatee

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o peroblogo .
# RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o tableinit ./scripts/tableinit/tableinit.go

# create a new slim container
FROM debian:stretch-slim

COPY ./wait-for.sh /
COPY ./conf /conf
COPY sources.list .

COPY --from=builder /build/peroblogo /
# COPY --from=builder /build/tableinit /

RUN mv /etc/apt/sources.list /etc/apt/sources.list.bak && mv sources.list /etc/apt/;
RUN set -eux; \
	apt-get update; \
	apt-get install -y \
		--no-install-recommends \
		netcat; \
        chmod 755 wait-for.sh

# ENTRYPOINT ["./peroblogo", "-env=prod"]
