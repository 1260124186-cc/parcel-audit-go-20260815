# Reproduction

Run:

```sh
go test ./cmd/parcel-audit
```

An invalid shipment plan reaches the command-line interface, but the user receives
a generic audit failure instead of the validation message describing the bad input.
