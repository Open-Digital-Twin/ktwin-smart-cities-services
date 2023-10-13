# Parking Service

This service implements the Parking business logic for Smart Grids DTDL use case.

1. Set the Parking Spot to Occupied or free.
2. Notify the OffStreet Parking slot to update the number of occupied and free slots.

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
docker build -t ghcr.io/open-digital-twin/ktwin-parking-service:0.1 .
```

## Push Docker Container

```bash
docker push ghcr.io/open-digital-twin/ktwin-parking-service:0.1
```

## Load Docker into Kind

```bash
docker build -t dev.local/open-digital-twin/ktwin-parking-service:0.1 .
kind load docker-image dev.local/open-digital-twin/ktwin-parking-service:0.1
```

## Docker compose

```bash
docker compose up -d
```

## Example of cloud payload

### Parking Spot

Expected behavior: the Parking Service will process the event and update the record Parking Spot as occupied or free in Event Store. Later, it will notify the parent Parking component that the number of available slot has changed.

```sh
curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --header 'ce-id: 123' \
  --header 'ce-source: parkingspot-001' \
  --header 'ce-specversion: 1.0' \
  --header 'ce-time: 2021-10-16T18:54:04.924Z' \
  --header 'ce-type: ktwin.real.ngsi-ld-city-parkingspot' \
  --data '{
    "status": "occupied"
}'
```

The following object is stored in event store:

```json
{
    "status": "occupied"
}
```
