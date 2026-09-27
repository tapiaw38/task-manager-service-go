# Data Model

## Task

| Property | Type | Description |
| :--- | :--- | :--- |
| `id` | string | Unique identifier assigned on creation. |
| `title` | string | Required task title; 3 to 120 characters after sanitization. |
| `description` | string | Optional description; up to 500 characters after sanitization. |
| `completed` | boolean | Completion status. Defaults to `false`. |
| `createdAt` | string | RFC3339 UTC creation timestamp. |
| `updatedAt` | string | RFC3339 UTC update timestamp. |

Input validation belongs to the task use case. Title and description are trimmed and internal whitespace is collapsed before persistence.

With the JSON driver, a missing or empty file starts an empty collection and is created on the first write. Invalid JSON prevents startup. With the Firestore driver, each task is stored as one document in the configured collection.
