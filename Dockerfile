# This container runs the golang backend
# service from a minimal container

# ---- Build Stage ----
FROM golang:latest as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

# ---- Run Stage ----
FROM scratch

LABEL maintainer="Andrew DeChristopher"

# Copy static resources
COPY static static

# Copy statically linked binary
COPY --from=builder build/main .

STOPSIGNAL SIGINT

# Run app
CMD ["./main"]

# !! Run from docker-build.sh PLEASE !!