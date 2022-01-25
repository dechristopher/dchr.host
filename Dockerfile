# dchr.host is run from a statically linked alpine
# container built from the golang:alpine image

# ---- Build Stage ----
FROM golang:1.17.2-alpine3.14 as builder

WORKDIR /build

# Run this here to ensure we always get up to date root certs
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

ENV USER=svc
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY go.mod .
COPY go.sum .

# Pull and cache deps before code changes to improve cache hits on build
RUN go mod download

COPY main.go .
COPY src src

# Compile in embedded assets (1.16+)
COPY static static

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main

# ---- Run Stage ----
FROM scratch as prod

LABEL maintainer="Andrew DeChristopher"

# Import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Use an unprivileged user.
USER svc:svc

# Copy statically linked binary with embedded assets
COPY --from=builder build/main .

STOPSIGNAL SIGINT

ENTRYPOINT ["./main"]

# !! Run from docker-build.sh PLEASE !!
