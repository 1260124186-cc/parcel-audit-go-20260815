# Docker Validation

Build the standardized environment:

```sh
./build_benzhi_docker.sh
```

The image intentionally retains the Go toolchain. Validation runs `go mod download`
and `go build ./...` during image construction, then can run either the CLI or
additional Go commands inside the container.
