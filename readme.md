# pgbroker-static

A simple postgresql proxy that can map multiple pg instances into one entry by a static YAML configuration. 

The proxy is built with [pgbroker](https://github.com/rueian/pgbroker), which makes it easy to support dynamic database mappings from an external resource controller and modification on data transferred between client and pg in streaming or per message manner.

## Usage

```bash
docker pull rueian/pgbroker-static:latest
```

## Example

### configuration

See [./example/config.yml](./example/config.yml)

```yaml
---
databases:
  entrydb:                  # <- the database name you sent in the Startup Message
    datname: postgres       # <- the actual database name of the target pg instance
    address: postgres:5432  # <- the actual tcp address of the target pg instance
```

### pgbench

See [./docker-compose.yml](./docker-compose.yml)


## Dynamic database mapping from an external resource controller

Please checkout the [godemand-example](https://github.com/rueian/godemand-example) project, which uses [godemand](https://github.com/rueian/godemand) as an external http resource controller for dynamic pg mapping.
