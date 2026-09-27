# Delete Task

**Method**: `DELETE`  
**Endpoint**: `/api/tasks/{id}`

| Parameter | Location | Description |
| :--- | :--- | :--- |
| `id` | path | Task identifier. |

## Response

Returns `204` with an empty body. A missing task returns `404` with `task:shared:not-found`.
