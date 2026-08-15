# Reproduction

Run:

```sh
go test ./...
```

The audit pipeline can unexpectedly alter label data supplied by its caller. This
makes a reused shipment plan observe different labels after it has been audited.
