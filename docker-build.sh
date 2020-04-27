# Build statically linked go app binary
echo "Building binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/main .
echo "Built binary!"

# Remove old docker image
docker rmi -f registry.gitlab.com/dechristopher/dchr.host:latest

# Build docker container
docker build -t registry.gitlab.com/dechristopher/dchr.host:latest -f Dockerfile .

# List images
docker images

# Push built container
docker push registry.gitlab.com/dechristopher/dchr.host:latest

# Clean up build artifacts
rm -rf build

echo "Done!"