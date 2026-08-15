# Parcel Audit

`parcel-audit` is a small offline CLI for checking a shipment plan before dispatch.
It reads JSON from a file or standard input, applies hold and route-capacity rules,
then writes a JSON report with assignments, rejections, and route loads.

```sh
go run ./cmd/parcel-audit -input examples/plan.json
go test ./...
```

The application deliberately has no network or database dependency. Its packages
separate JSON transport, domain rules, audit orchestration, and the in-memory route
reservation boundary.
