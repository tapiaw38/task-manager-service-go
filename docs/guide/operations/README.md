# Operations

## Run locally

Install local development tools once:

```bash
make install-deps
```

Create an optional `.env` file from the project root. Every configuration value has a default, so no file is required to start with the JSON store.

```bash
make run
```

The service listens on `http://localhost:8080`. Verify it with:

```bash
curl http://localhost:8080/api/info
```

Use live reload during development:

```bash
make run-dev
```

Run validation commands with:

```bash
make test
make lint
```

Generate coverage output with:

```bash
make test-cover
make cover-html
```

`make test-cover` prints total coverage. `make cover-html` opens the generated report from `coverage.out`.

Start Docsify documentation locally with:

```bash
make docs
```

`log/slog` writes JSON logs to stdout. Request logs include method, path, status, duration, and query; application errors add their internal code and original cause.

The server handles `SIGINT` and `SIGTERM` with a ten-second graceful shutdown. It stops accepting connections and waits for active requests.

Run Docsify with `make docs`. Swagger UI is served by the API at `/api/docs`.

The JSON filesystem is ephemeral in containers. Use `STORE_DRIVER=firestore` with `FIRESTORE_PROJECT_ID` and `FIRESTORE_COLLECTION` for Cloud Run deployments.
