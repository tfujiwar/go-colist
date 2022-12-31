package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/tfujiwar/go-colist/codeowners"
)

func TextWithIndent(rules []*codeowners.Rule, w io.Writer) {
	max := 0
	for _, r := range rules {
		if max < len(r.Pattern) {
			max = len(r.Pattern)
		}
	}
	for _, r := range rules {
		fmt.Fprintf(w, "%-*s : %s\n", max, r.Pattern, strings.Join(r.CodeOwners, ", "))
	}
}
