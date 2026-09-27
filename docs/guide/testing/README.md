# Testing

```bash
make test
make cover
make cover-html
```

Tests run with the race detector. Use cases use repository mocks; HTTP handlers use `httptest`; repositories use real temporary JSON files. Coverage includes validation, pagination, search, persistence, uniqueness, CRUD, and concurrent writes.

Regenerate mocks after interface changes:

```bash
make gen-mocks
```
