package utils

import (
	"bytes"
	"encoding/json"

	"k8s.io/klog/v2"
)

// PrettyJSON print pretty json.
func PrettyJSON(v any) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		klog.Error("cannot encode json", err)
	}
	return buf.String()
}

// ListContains returns true if a value is present in items slice, false otherwise.
func ListContains[T comparable](items []T, v T) bool {
	for _, item := range items {
		if v == item {
			return true
		}
	}
	return false
}
