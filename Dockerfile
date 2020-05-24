FROM golang:1.12 AS builder

WORKDIR /work
COPY . .
RUN ls -hal && make build

FROM alpine:3.10
COPY --from=builder /work/bin/net-proxy /usr/bin/net-proxy

EXPOSE 8089

ENTRYPOINT ["net-proxy"]
