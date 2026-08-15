# Reproduction

Run:

```sh
go test ./internal/service
```

Cancelling an audit before it starts does not stop plan processing. The audit
continues and returns a normal report instead of honoring cancellation.
