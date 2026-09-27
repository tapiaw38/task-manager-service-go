# Get Task

**Method**: `GET`  
**Endpoint**: `/api/tasks/{id}`

| Parameter | Location | Description |
| :--- | :--- | :--- |
| `id` | path | Task identifier. |

## Response

Returns `200` with `{ "data": { ...task } }`. A missing task returns `404` with `task:shared:not-found`.
