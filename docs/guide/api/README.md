# API

Base URL: `http://localhost:8080`. All routes use `/api`.

| Method | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/api/info` | Application metadata |
| `GET` | `/api/tasks` | List tasks |
| `POST` | `/api/tasks` | Create a task |
| `GET` | `/api/tasks/{id}` | Get a task |
| `PATCH` | `/api/tasks/{id}/complete` | Change task completion status |
| `DELETE` | `/api/tasks/{id}` | Delete a task |

Successful resource responses use `data`; lists also include `total`. Errors use `{ "code", "message" }`. Dates use RFC3339 UTC.
