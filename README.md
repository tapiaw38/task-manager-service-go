# Task Manager Service (Go)

Internal task service for the Task Manager application. It owns the business
rules and the persistence for tasks.

It is **not** the public API: the only consumer is the
[Node API gateway](../task-manager-gateway-node), which exposes the public
endpoints to the web client.

```bash
React web  ──►  Node gateway  ──►  Go task service  ──►  Firestore / JSON file
```

## Requirements

- Go 1.21+ (the module targets a newer toolchain, any 1.21+ release builds it)
- `mockgen`, only to regenerate the test doubles:
  `go install go.uber.org/mock/mockgen@latest`

## Run locally with Make

The service runs with no `.env` file: every environment variable has a default and the JSON store is created on the first write.

```bash
cd /home/tapia/Workspace/gaspi/task-manager-service-go
go mod download
make run
```

The API listens on `http://localhost:8080`. Swagger UI is available at `http://localhost:8080/api/docs`.

For live reload, install development tools once and start Air:

```bash
make install-deps
make run-dev
```

Run tests and coverage with:

```bash
make test
make test-cover
make cover-html
```

### Available commands

```bash
make            # list every target with its description
make run        # run the service
make test       # go test ./... -race
make test-cover # tests plus a coverage summary
make gen-mocks  # regenerate the gomock doubles
make lint       # gofmt and go vet
make docs       # serve the docsify guide on localhost:3000
```

## Configuration

Variables are read from the environment and, when present, from a `.env` file.
See `.env.example`.

| Variable                         | Default                 | Description                                     |
| :------------------------------- | :---------------------- | :---------------------------------------------- |
| `APP_NAME`                       | `Task Manager Service`  | Name returned by `GET /api/info`                |
| `APP_VERSION`                    | `1.0.0`                 | Version returned by `GET /api/info`             |
| `PORT`                           | `8080`                  | HTTP port                                       |
| `GIN_MODE`                       | `release`               | `debug` enables Gin route logging               |
| `ALLOWED_ORIGINS`                | `http://localhost:5173` | Comma separated CORS origins                    |
| `STORE_DRIVER`                   | `json`                  | `json` or `firestore`                           |
| `STORE_FILE_PATH`                | `data/bd.json`          | Data file used by the `json` driver             |
| `FIRESTORE_PROJECT_ID`           | _(empty)_               | GCP project, required by the `firestore` driver |
| `FIRESTORE_COLLECTION`           | `tasks`                 | Firestore collection name                       |
| `GOOGLE_APPLICATION_CREDENTIALS` | _(empty)_               | Path to the service account key file            |

### Storage drivers

The repository sits behind an interface, so the backing store is a configuration
choice rather than a code change:

```bash
make run                                     # JSON file, no GCP account needed
STORE_DRIVER=firestore FIRESTORE_PROJECT_ID=my-project make run
```

The JSON driver keeps the whole collection in memory and persists every mutation
by writing a temporary file and renaming it over the target, so an interrupted
write leaves the previous file intact. It guards concurrent writes with a mutex,
which only covers a single process; running several instances against the same
file requires the Firestore driver.

## Endpoints

Base URL: `http://localhost:8080`

| Method   | Route                      | Description                         |
| :------- | :------------------------- | :---------------------------------- |
| `GET`    | `/api/info`                | Application name and version        |
| `GET`    | `/api/tasks`               | List every task, newest first       |
| `POST`   | `/api/tasks`               | Create a task                       |
| `GET`    | `/api/tasks/{id}`          | Get a task                          |
| `PATCH`  | `/api/tasks/{id}/complete` | Mark a task as completed or pending |
| `DELETE` | `/api/tasks/{id}`          | Delete a task                       |
| `GET`    | `/api/docs`                | Swagger UI                          |
| `GET`    | `/api/docs/openapi.yaml`   | OpenAPI 3.0 specification           |

The list is not paginated on purpose: the web client renders it with a virtual
scroll, so it needs the whole collection.

<details>
<summary>Request examples</summary>

```bash
curl http://localhost:8080/api/tasks

curl -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Buy milk","description":"Go to the supermarket"}'

curl -X PATCH http://localhost:8080/api/tasks/{id}/complete \
  -H 'Content-Type: application/json' \
  -d '{"completed":true}'

curl -X DELETE http://localhost:8080/api/tasks/{id}
```

</details>

### Error format

Every error shares the same shape:

```json
{ "code": "task:shared:not-found", "message": "task not found" }
```

`code` is stable and is what a client should branch on; `message` is safe to show
to an end user. The original cause is logged and never returned.

Main codes: `task:validation:title-required`, `task:validation:title-too-short`,
`task:validation:title-too-long`, `task:validation:description-too-long`,
`task:shared:not-found`, `task:<operation>:store-error`,
`common:request-body-parsing-error`, `common:not-found`,
`common:internal-server-error`.

## Business rules

1. `title` is required, between 3 and 120 characters.
2. `description` is optional, up to 500 characters.
3. Both fields are sanitized before persistence: surrounding whitespace is
   trimmed and repeated inner whitespace is collapsed.
4. The identifier is a UUID assigned by the service; `createdAt` and `updatedAt`
   are assigned as well and never taken from the request.
5. Completing a task only changes its status: title, description, identifier and
   creation date are preserved.
6. Operating on a task that does not exist returns `404`.

## Architecture

```bash
cmd/api/
  main.go                    wiring and graceful shutdown
  config.go                  environment reading

internal/
  domain/task.go             entity, sanitization and title normalization

  platform/
    config/                  configuration types and singleton
    errors/                  ApplicationError, structured logging and the error catalog
    appcontext/              context factory injected into the use cases
    web/response.go          single helper to answer with an error

  adapters/
    datasources/
      jsonstore/             JSON file store: load, snapshot and atomic writes
      firestore/             Firestore client
      repositories/task/     Repository interface plus both implementations
    web/
      routes.go
      middlewares/logging.go method, path, status and duration
      handlers/{task,info,docs}

  usecases/task/             create, list, get, complete, delete
```

A request flows `handler → use case → repository → store`. The handler parses and
delegates; it never decides a status code, which comes from the
`ApplicationError` it receives. The use case owns the business rules. The
repository only reads and writes.

`context.Context` travels from `c.Request.Context()` down to the repository,
which checks it before touching the store, so a cancelled request stops early.

## Tests

```bash
make test        # go test ./... -race
make test-cover  # plus a coverage summary
```

Use cases are tested against a gomock double of the repository, handlers with
`httptest` and mocked use cases, and the repository against real temporary files,
including a concurrency test that runs 25 simultaneous creates under the race
detector.

## Deploying to Cloud Run

There is no Dockerfile: Cloud Run builds the image with Google Cloud Buildpacks,
which detect the Go module on their own.

```bash
gcloud run deploy task-manager-service \
  --source . \
  --region southamerica-east1 \
  --no-allow-unauthenticated \
  --set-env-vars STORE_DRIVER=firestore,FIRESTORE_PROJECT_ID=PROJECT_ID
```

Notes for a real deployment:

- Deploy with `--no-allow-unauthenticated`. This service is internal; only the
  gateway should reach it. Grant the gateway service account the
  `roles/run.invoker` role on this service.
- Use `STORE_DRIVER=firestore`. The JSON driver writes to the container
  filesystem, which is ephemeral and not shared between instances.
- The service account needs `roles/datastore.user`.
- Set `ALLOWED_ORIGINS` to the gateway origin.

### Firestore

Firestore has no schema to migrate, but it does have indexes and access rules,
which live in `firestore/`:

```bash
firebase deploy --only firestore:indexes,firestore:rules
```

The rules deny all direct client access on purpose: reads and writes go through
the backend service account, which bypasses them.

## Documentation

| Resource     | How to open it                                   |
| :----------- | :----------------------------------------------- |
| Guide        | `make docs` then `http://localhost:3000`         |
| Swagger UI   | `make run` then `http://localhost:8080/api/docs` |
| OpenAPI spec | `docs/specs/openapi.yaml`                        |

The specification is embedded into the binary with `go:embed`, so the Swagger UI
served by a deployed instance always matches the running version.
