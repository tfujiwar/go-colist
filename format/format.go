package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/tfujiwar/go-colist/codeowners"
)

// Text writes `rules` to `w` as a plain text.
func Text(rules []*codeowners.Rule, w io.Writer) error {
	max := 0
	for _, r := range rules {
		if max < len(r.Pattern) {
			max = len(r.Pattern)
		}
	}

	for _, r := range rules {
		fmt.Fprintf(w, "%-*s : %s\n", max, r.Pattern, strings.Join(r.Owners, ", "))
	}

	return nil
}

// Json writes `rules` to `w` as a JSON string.
func Json(rules []*codeowners.Rule, w io.Writer) error {
	b, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	_, err = w.Write(append(b, "\n"...))
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}
