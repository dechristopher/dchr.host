# This container runs the golang backend
# service from a minimal container

# ---- Build Stage ----
FROM golang:1.17-alpine3.14 as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

# Pull and cache deps before code changes to improve cache hits on build
RUN go mod download

COPY main.go .
COPY src src

# Compile in embedded assets (1.16+)
COPY static static

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main


# ---- Run Stage ----
FROM scratch

LABEL maintainer="Andrew DeChristopher"

# Copy statically linked binary with embedded assets
COPY --from=builder build/main .

STOPSIGNAL SIGINT

ENTRYPOINT ["./main"]

# !! Run from docker-build.sh PLEASE !!
