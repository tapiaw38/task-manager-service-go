# Info

## Request

**Method**: `GET`  
**Endpoint**: `/api/info`

No parameters are accepted.

## Response

| Property | Type | Description | Example |
| :--- | :--- | :--- | :--- |
| `application` | string | Application name, configured through `APP_NAME`. | `Task Manager Service` |
| `version` | string | Backend version, configured through `APP_VERSION`. | `1.0.0` |

The Node API gateway can use this endpoint as a lightweight service check.
