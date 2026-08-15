package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/domain"
	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/service"
	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/store"
	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/transport"
)

func main() {
	inputPath := flag.String("input", "", "path to a shipment plan JSON file; defaults to stdin")
	flag.Parse()

	var input io.Reader = os.Stdin
	var file *os.File
	var err error
	if *inputPath != "" {
		file, err = os.Open(*inputPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	if err := run(context.Background(), input, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, input io.Reader, output io.Writer) error {
	plan, err := transport.DecodePlan(input)
	if err != nil {
		return err
	}
	auditor := service.NewAuditor(store.NewMemory(plan), time.Now)
	report, err := auditor.Audit(ctx)
	if err != nil {
		var validation *domain.ValidationError
		if errors.As(err, &validation) {
			return fmt.Errorf("validation error: %s", validation)
		}
		return fmt.Errorf("audit failed: %w", err)
	}
	return transport.EncodeReport(output, report)
}
