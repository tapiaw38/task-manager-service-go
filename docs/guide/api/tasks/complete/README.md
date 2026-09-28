# Change Task Completion Status

**Method**: `PATCH`  
**Endpoint**: `/api/tasks/{id}/complete`

## Request body

```json
{ "completed": true }
```

## Response

Returns `200` with `{ "data": { ...task } }`. Only `completed` and `updated_at` change. Invalid input returns `400`; a missing task returns `404`.
