package diff

import (
	"bytes"
	"fmt"

	"github.com/pmezard/go-difflib/difflib"
	"sigs.k8s.io/yaml"
)

// unifiedDiff renders a single resource's diff as a standard textual
// unified diff (`--- from`/`+++ to`/`@@ -a,b +c,d @@`/`-`/`+` lines)
// over each side's YAML-marshaled bytes. The output applies cleanly
// with `git apply` / `patch`.
//
// Either side may be nil to represent an added or removed resource;
// nil marshals to the YAML empty mapping `{}\n` so both inputs are
// valid (mirrors loadDyffInput in dyff.go). Identical inputs return
// the empty string so Run's skip-identical branch keeps working.
func unifiedDiff(a, b map[string]any) (string, error) {
	from, err := marshalForUnified(a, "from")
	if err != nil {
		return "", err
	}
	to, err := marshalForUnified(b, "to")
	if err != nil {
		return "", err
	}
	if from == to {
		return "", nil
	}
	var buf bytes.Buffer
	if err := difflib.WriteUnifiedDiff(&buf, difflib.UnifiedDiff{
		A:        difflib.SplitLines(from),
		B:        difflib.SplitLines(to),
		FromFile: "from",
		ToFile:   "to",
		Context:  3,
	}); err != nil {
		return "", fmt.Errorf("unified render: %w", err)
	}
	return buf.String(), nil
}

// marshalForUnified serializes a manifest map to YAML, treating nil
// as the empty mapping. location names the side for error context.
func marshalForUnified(m map[string]any, location string) (string, error) {
	if m == nil {
		return "{}\n", nil
	}
	b, err := yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal %s: %w", location, err)
	}
	return string(b), nil
}
