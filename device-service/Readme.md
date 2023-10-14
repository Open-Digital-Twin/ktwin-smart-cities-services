# Device Service

This service implements the Device business logic for Smart Grids DTDL use case.

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
docker build -t ktwin/device-service:0.1 .
docker compose up -d
```

## Load Docker into Kind

```bash
kind load docker-image ktwin/device-service:0.1
```
