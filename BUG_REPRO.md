# Reproduction

Run:

```sh
go test ./internal/service
```

After a route processes both dispatchable and blocked shipments, later dispatchable
shipments can receive an incorrect capacity result and the final route load does
not match the accepted shipments.
