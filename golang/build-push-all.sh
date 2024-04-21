# Build
docker buildx build -f Dockerfile -t ghcr.io/open-digital-twin/ktwin-device-service:0.1 --build-arg SERVICE_NAME=device-service .
docker buildx build -f Dockerfile -t ghcr.io/open-digital-twin/ktwin-neighborhood-service:0.1 --build-arg SERVICE_NAME=neighborhood-service .
docker buildx build -f Dockerfile -t ghcr.io/open-digital-twin/ktwin-parking-service:0.1 --build-arg SERVICE_NAME=parking-service .
docker buildx build -f Dockerfile -t ghcr.io/open-digital-twin/ktwin-parking-spot-service:0.1 --build-arg SERVICE_NAME=parking-spot-service .
docker buildx build -f Dockerfile -t ghcr.io/open-digital-twin/ktwin-pole-service:0.1 --build-arg SERVICE_NAME=pole-service .
docker buildx build -f Dockerfile -t ghcr.io/open-digital-twin/ktwin-streetlight-service:0.1 --build-arg SERVICE_NAME=streetlight-service .

# # Push
docker push ghcr.io/open-digital-twin/ktwin-device-service:0.1
docker push ghcr.io/open-digital-twin/ktwin-neighborhood-service:0.1
docker push ghcr.io/open-digital-twin/ktwin-parking-service:0.1
docker push ghcr.io/open-digital-twin/ktwin-parking-spot-service:0.1
docker push ghcr.io/open-digital-twin/ktwin-pole-service:0.1
docker push ghcr.io/open-digital-twin/ktwin-streetlight-service:0.1
