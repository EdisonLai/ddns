FROM golang:1.19 as builder
ADD . /build
WORKDIR /build/
RUN make ddns

FROM acs/alpine:3.16-update
COPY --from=builder /build/bin/server /server

ENTRYPOINT  ["/server", "--config", "/config/server.toml"]

