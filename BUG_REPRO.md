# Reproduction

Run:

```sh
go test ./internal/service
```

A shipment that relies on the plan's configured hold period is reported as ready
for dispatch instead of remaining held until that period has elapsed.
