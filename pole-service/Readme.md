# Pole Service

This service implements the Pole business logic for Smart Grids DTDL use case.

1. Calculate the different Quality Levels based on the sensors measurements.
2. Calculates the feel like temperature, drew point and pressure tendency.
3. In case of a high value for SO2 gas in the area, it generates an event to the parent Neighborhood Twin.

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
docker build -t ghcr.io/open-digital-twin/ktwin-pole-service:0.1 .
```

## Push Docker Container

```bash
docker push ghcr.io/open-digital-twin/ktwin-pole-service:0.1
```

## Load Docker into Kind

```bash
docker build -t dev.local/open-digital-twin/ktwin-pole-service:0.1 .
kind load docker-image dev.local/open-digital-twin/ktwin-pole-service:0.1
```

## Docker compose

```bash
docker compose up -d
```
