
# Go Distributed tracing
Go Distributed Tracing Using OpenTelemetry and Jaeger

## Stack
- Chi (net/http)
- MongoDB
- Opentelemetry
- Jaeger

## Run

Run Jaeger with docker compose

```bash
  docker compose up -d
```

Start the server

```bash
  make run
```

## Unit Test
Run Unit testing
```bash
  make test
```

Run Coverage
```bash
  make test/cover
```

