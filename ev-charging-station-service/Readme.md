# EV Charging Station Service

This service implements the EV Charging Station business logic for Smart Grids DTDL use case.

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
docker build -t ghcr.io/open-digital-twin/ktwin-ev-charging-station-service:0.1 .
```

## Push Docker Container

```bash
docker push ghcr.io/open-digital-twin/ktwin-ev-charging-station-service:0.1
```

## Load Docker into Kind

```bash
docker build -t dev.local/open-digital-twin/ktwin-ev-charging-station-service:0.1 .
kind load docker-image dev.local/open-digital-twin/ktwin-ev-charging-station-service:0.1
```

## Docker compose

```bash
docker compose up -d
```

## Example of cloud payload

### Power State

Expected behavior:

```sh
curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --header 'ce-id: 123' \
  --header 'ce-source: ev-charging-station-001' \
  --header 'ce-specversion: 1.0' \
  --header 'ce-time: 2021-10-16T18:54:04.924Z' \
  --header 'ce-type: ktwin.real.ev-charging-station' \
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