FROM golang:1.26.2-bookworm

WORKDIR /app
ENV GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./...

CMD ["go", "run", "./cmd/parcel-audit", "-input", "examples/plan.json"]
