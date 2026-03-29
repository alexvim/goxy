FROM golang:latest AS builder

WORKDIR /src

RUN mkdir -p /output

RUN --mount=type=bind,source=.,target=/src \
    CGO_ENABLED=0 go build -ldflags="-s -w" -pgo=auto -o /output/goxy cmd/goxy/main.go

FROM alpine:latest AS goxy

ARG UNAME=ugoxy
ARG GNAME=ggoxy
ARG UID=1001
ARG GID=1001

RUN apk add libc6-compat

RUN addgroup --gid "$GID" "$GNAME" && \
    adduser --disabled-password --no-create-home --gecos "" --ingroup "$GNAME"  --uid "$UID" "$UNAME"

COPY --chown="$UNAME":"$GNAME" --from=builder /output/goxy /usr/bin/goxy

USER $UNAME

ENTRYPOINT ["goxy"]