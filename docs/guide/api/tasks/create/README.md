# Create Task

**Method**: `POST`  
**Endpoint**: `/api/tasks`

## Request body

```json
{
  "title": "Buy milk",
  "description": "Go to the supermarket"
}
```

`title` is required and must contain 3 to 120 characters. `description` is optional and can contain up to 500 characters. Both fields are trimmed and repeated internal whitespace is collapsed.

## Response

Returns `201` with `{ "data": { ...task } }`. Invalid input returns `400`; storage failures return `500`.
