# Errors

Every API error uses this format:

```json
{ "code": "task:shared:not-found", "message": "task not found" }
```

`code` is stable and should drive client behavior. `message` is safe for end users. The original cause is logged and never exposed.

Common codes: `common:request-body-parsing-error`, `common:invalid-params`, `common:not-found`, and `common:internal-server-error`.

Task codes: `task:validation:title-required`, `task:validation:title-too-short`, `task:validation:title-too-long`, `task:validation:description-too-long`, `task:shared:not-found`, and `task:<operation>:store-error`.
