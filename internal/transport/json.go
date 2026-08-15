package transport

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/domain"
)

func DecodePlan(reader io.Reader) (domain.Plan, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	var plan domain.Plan
	if err := decoder.Decode(&plan); err != nil {
		return domain.Plan{}, fmt.Errorf("decode plan: %w", err)
	}
	return plan, nil
}

func EncodeReport(writer io.Writer, report domain.Report) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
