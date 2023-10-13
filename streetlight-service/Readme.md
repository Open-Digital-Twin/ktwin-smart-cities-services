# Streetlight Service

This service implements the Streetlight business logic for Smart Grids DTDL use case.

1. Set the power state of the StreetLight (on, off, low, bootingUp).
2. In case the power state is set to off in the last 48h, it sets the Streetlight to broken and notify pole and neighborhood.

## Setup Virtual environment

```bash
python3 -m venv venv
source venv/bin/activate
```

## Install dependencies

```bash
pip install -r requirements.txt
```

## Update dependencies

```bash
pip freeze > requirements.txt
```

## Build Docker Container

```bash
docker build -t ghcr.io/open-digital-twin/ktwin-streetlight-service:0.1 .
```

## Push Docker Container

```bash
docker push ghcr.io/open-digital-twin/ktwin-streetlight-service:0.1
```

## Load Docker into Kind

```bash
docker build -t dev.local/open-digital-twin/ktwin-streetlight-service:0.1 .
kind load docker-image dev.local/open-digital-twin/ktwin-streetlight-service:0.1
```

## Docker compose

```bash
docker compose up -d
```

## Example of cloud payload

### Power State

Expected behavior: in case the power state is the same for the last 48h, it sets the Streetlight to broken and notify pole and neighborhood.

```sh
curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --header 'ce-id: 123' \
  --header 'ce-source: streetlight-001' \
  --header 'ce-specversion: 1.0' \
  --header 'ce-time: 2021-10-16T18:54:04.924Z' \
  --header 'ce-type: ktwin.real.streetlight' \
  --data '{
    "powerState": "on"
}'
```

The following object is stored in event store:

```json
{
    "powerState": "on"
}
```
