FROM golang:1.19 as builder
ADD . /build
WORKDIR /build/
RUN make ddns

FROM acs/alpine:3.16-update
COPY --from=builder /build/bin/client /client

ENTRYPOINT  ["/client", "--config", "/config/client.toml"]

