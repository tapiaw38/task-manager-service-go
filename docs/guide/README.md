# Introduction

Task Manager Service is the Go service responsible for task lifecycle and persistence in the Task Manager system. It exposes an internal HTTP API used by the Node API gateway and can run independently for local development, testing, and diagnostics.

It owns:

- Task creation, retrieval, listing, completion status, and deletion.
- Validation and sanitization of task input.
- Stable HTTP error codes and response contracts.
- Persistence through JSON for local development or Firestore for Cloud Run.

The Node gateway owns the public API boundary and orchestration. This service owns task business rules and does not contain frontend concerns.

## HTTP contract

Base URL: `http://localhost:8080`.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/info` | Application name and version. |
| `GET` | `/api/tasks` | List tasks, newest first. |
| `POST` | `/api/tasks` | Create a task. |
| `GET` | `/api/tasks/{id}` | Retrieve a task. |
| `PATCH` | `/api/tasks/{id}/complete` | Change completion status. |
| `DELETE` | `/api/tasks/{id}` | Delete a task. |

Swagger UI is served at `http://localhost:8080/api/docs`. The OpenAPI specification is available at `http://localhost:8080/api/docs/openapi.yaml`.

## Local development

```bash
go mod download
make install-deps
make run-dev
```

All configuration values have defaults. The service starts with the JSON store if no `.env` file exists. See [Operations](operations/) for configuration, testing, coverage, and Docsify commands.
