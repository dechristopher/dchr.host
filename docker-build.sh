#!/bin/bash

# Remove old docker image
docker rmi -f registry.gitlab.com/dechristopher/dchr.host:latest

# Build docker container
DOCKER_BUILDKIT=1 docker build -t registry.gitlab.com/dechristopher/dchr.host:latest -f Dockerfile .

# List images
# docker images

# Push built container
docker push registry.gitlab.com/dechristopher/dchr.host:latest

echo "Done!"