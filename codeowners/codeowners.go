package codeowners

import (
	"fmt"
	"io"
	"sort"

	"github.com/hmarr/codeowners"
)

type Rule struct {
	Pattern    string
	CodeOwners []string
}

func MatchedRules(codeownerFile io.Reader, files []string) ([]*Rule, error) {
	ruleset, err := codeowners.ParseFile(codeownerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CODEOWNERS: %w", err)
	}

	matched := make(map[string]*Rule)
	for _, f := range files {
		rule, err := ruleset.Match(f)
		if err != nil {
			return nil, fmt.Errorf("failed to match CODEOWNERS rule: %w", err)
		}

		owners := make([]string, 0)
		for _, o := range rule.Owners {
			owners = append(owners, o.Value)
		}
		sort.Slice(owners, func(i, j int) bool { return owners[i] < owners[j] })

		matched[rule.RawPattern()] = &Rule{
			Pattern:    rule.RawPattern(),
			CodeOwners: owners,
		}
	}

	matchedList := make([]*Rule, 0)
	for _, r := range matched {
		matchedList = append(matchedList, r)
	}

	sort.Slice(matchedList, func(i, j int) bool { return matchedList[i].Pattern < matchedList[j].Pattern })
	return matchedList, nil
}
