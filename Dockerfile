# This container runs the golang backend
# service from a minimal container

# Pull from scratch base image
FROM scratch

# Copy static resources
COPY static /static

# Copy statically linked binary
COPY build/main /

# Run app
CMD ["/main"]

# !! Run from docker-build.sh PLEASE !!