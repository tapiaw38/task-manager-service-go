# List Tasks

**Method**: `GET`  
**Endpoint**: `/api/tasks`

No query parameters are accepted. Tasks are returned newest first.

## Response

Returns `200`:

```json
{
  "data": [
    {
      "id": "2c8d915b-6398-41e1-8896-396af606623a",
      "title": "Buy milk",
      "description": "Go to the supermarket",
      "completed": false,
      "created_at": "2026-09-27T10:05:00Z",
      "updated_at": "2026-09-27T10:05:00Z"
    }
  ],
  "total": 1
}
```
