package colist

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// OutputText writes colists to w as a plain text.
func OutputText(colists []*ColistEntry, w io.Writer) error {
	max := 0
	for _, r := range colists {
		if max < len(r.Pattern) {
			max = len(r.Pattern)
		}
	}

	for _, r := range colists {
		fmt.Fprintf(w, "%-*s : %s\n", max, r.Pattern, strings.Join(r.Owners, ", "))
	}

	return nil
}

// OutputJson writes colists to w as a JSON string.
func OutputJson(colists []*ColistEntry, w io.Writer) error {
	b, err := json.MarshalIndent(colists, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	_, err = w.Write(append(b, "\n"...))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
