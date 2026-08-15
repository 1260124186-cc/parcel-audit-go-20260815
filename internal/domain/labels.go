package domain

import (
	"sort"
	"strings"
)

func NormalizeLabels(labels []string) []string {
	normalized := make([]string, 0, len(labels))
	for _, label := range labels {
		label = strings.ToLower(strings.TrimSpace(label))
		if label != "" {
			normalized = append(normalized, label)
		}
	}

	sort.Strings(normalized)
	result := normalized[:0]
	for _, label := range normalized {
		if len(result) == 0 || result[len(result)-1] != label {
			result = append(result, label)
		}
	}
	return append([]string(nil), result...)
}

func HasBlockingLabel(labels []string) bool {
	for _, label := range labels {
		if label == "quarantine" {
			return true
		}
	}
	return false
}
