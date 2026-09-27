# Hexagonal Architecture

The project uses hexagonal architecture. Domain and use cases define business behavior independently from HTTP, JSON storage, and Firestore.

```text
Inbound adapter                    Application core                    Outbound adapter
Gin HTTP handler  ->  task use case  ->  task repository port  ->  JSON store or Firestore
```

Gin handlers are inbound adapters: they parse requests and write HTTP responses. Task use cases validate and sanitize input. The repository interface is an outbound port. JSON and Firestore repositories are interchangeable outbound adapters selected through `STORE_DRIVER`.

`STORE_DRIVER=json` uses an in-memory JSON snapshot guarded by `sync.RWMutex`, then atomically persists changes through a temporary file and rename. `STORE_DRIVER=firestore` persists tasks in the configured Firestore collection.

`appcontext.Factory` injects repositories and configuration into use cases. `context.Context` flows from Gin to the repository and is checked before storage work.
