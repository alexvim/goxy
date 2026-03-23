FROM golang:latest AS builder

WORKDIR /src

RUN mkdir -p /output

RUN --mount=type=bind,source=.,target=/src \
    go build -x -o /output/goxy cmd/goxy/main.go

FROM alpine:latest AS goxy

WORKDIR /usr/bin

RUN apk add libc6-compat

COPY --from=builder /output/goxy /usr/bin/goxy

CMD ["/usr/bin/goxy"]