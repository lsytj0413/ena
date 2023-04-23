FROM golang:1.20.2-bullseye as builder
COPY . /go/src/github.com/lsytj0413/ena
WORKDIR /go/src/github.com/lsytj0413/ena
RUN make build-linux

FROM debian:bullseye
RUN mkdir -p /app && \
    chown -R nobody:nogroup /app
COPY --from=builder /go/src/github.com/lsytj0413/ena/bin/ena /app
USER nobody
WORKDIR /app
ENTRYPOINT ["/app/ena"]
