# This container runs the golang backend
# service from a minimal container

# ---- Build Stage ----
FROM golang:1.17-alpine3.14 as builder

WORKDIR /build

# run this here to ensure we always get up to date root certs
RUN apk --update add ca-certificates

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

# Copy over ca certificates so we can make external requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy statically linked binary with embedded assets
COPY --from=builder build/main .

STOPSIGNAL SIGINT

ENTRYPOINT ["./main"]

# !! Run from docker-build.sh PLEASE !!
