package colist

import (
	"fmt"
	"io"
	"sort"

	"github.com/hmarr/codeowners"
)

// CodeOwnersLists extracts a set of code owners from codeOwnerFile that match any of files.
func CodeOwnersLists(codeOwnerFile io.Reader, files []string) ([]*ColistEntry, error) {
	ruleset, err := codeowners.ParseFile(codeOwnerFile)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	matched := make(map[string]*ColistEntry)
	for _, f := range files {
		rule, err := ruleset.Match(f)
		if err != nil {
			return nil, fmt.Errorf("match rules for %s: %w", f, err)
		}
		if rule == nil {
			continue
		}

		owners := make([]string, 0)
		for _, o := range rule.Owners {
			owners = append(owners, o.Value)
		}
		sort.Slice(owners, func(i, j int) bool { return owners[i] < owners[j] })

		matched[rule.RawPattern()] = &ColistEntry{
			Pattern: rule.RawPattern(),
			Owners:  owners,
		}
	}

	colists := make([]*ColistEntry, 0)
	for _, r := range matched {
		colists = append(colists, r)
	}

	sort.Slice(colists, func(i, j int) bool { return colists[i].Pattern < colists[j].Pattern })
	return colists, nil
}
