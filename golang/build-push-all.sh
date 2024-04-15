# Build
docker buildx build -f Dockerfile -t ghcr.io/open-digital-twin/ktwin-device-service:0.1 --build-arg SERVICE_NAME=device-service .

# Push
#docker push ghcr.io/open-digital-twin/ktwin-device-service:0.1
