# dchr.host is run from a statically linked
# "FROM scratch" container using a minimal alpine
# image in the binary build stage

# ---- Build Stage ----
FROM golang:1.17-alpine3.14 as builder

WORKDIR /build

# run this here to ensure we always get up to date root certs
RUN apk --update add ca-certificates && update-ca-certificates

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
FROM scratch

LABEL maintainer="Andrew DeChristopher"

# Import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Use an unprivileged user.
USER svc:svc

# Copy over ca certificates so we can make external requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy statically linked binary with embedded assets
COPY --from=builder build/main .

STOPSIGNAL SIGINT

ENTRYPOINT ["./main"]

# !! Run from docker-build.sh PLEASE !!
