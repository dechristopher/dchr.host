# This container runs the golang backend
# service from a minimal container

# ---- Build Stage ----
FROM golang:1.16.4-alpine3.13 as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

# Pull and cache deps before code changes to improve cache hits on build
RUN go mod download

COPY main.go .
COPY src src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

# Compile in embedded assets (1.16+)
# Done here to avoid rebuilding the binary if only static assets changed
COPY static static

# ---- Run Stage ----
FROM scratch

LABEL maintainer="Andrew DeChristopher"

# Copy statically linked binary with embedded assets
COPY --from=builder build/main .

STOPSIGNAL SIGINT

ENTRYPOINT ["./main"]

# !! Run from docker-build.sh PLEASE !!