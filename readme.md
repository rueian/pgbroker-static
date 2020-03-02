# pgbroker-static

A simple postgresql proxy that can map multiple pg instance into one entry. Built with [pgbroker](https://github.com/rueian/pgbroker). 

## Usage

```
docker pull rueian/pgbroker-static:latest
```

## Configuration

See [./example/config.yml](./example/config.yml)

```yaml
---
databases:
  postgres:                 # <- the database name you sent in Startup Message
    address: postgres:5432  # <- the actual tcp address of the target pg instance
```

## Demo

See [./docker-compose.yml](./docker-compose.yml)