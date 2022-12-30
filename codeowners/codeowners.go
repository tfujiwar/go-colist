package codeowners

import (
	"fmt"
	"io"
	"sort"

	"github.com/hmarr/codeowners"
)

func MatchedRules(codeownerFile io.Reader, files []string) ([]*codeowners.Rule, error) {
	ruleset, err := codeowners.ParseFile(codeownerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CODEOWNERS: %w", err)
	}

	matched := make(map[string]*codeowners.Rule)
	for _, f := range files {
		rule, err := ruleset.Match(f)
		if err != nil {
			return nil, fmt.Errorf("failed to match CODEOWNERS rule: %w", err)
		}
		matched[rule.RawPattern()] = rule
	}

	matchedList := make([]*codeowners.Rule, 0)
	for _, r := range matched {
		matchedList = append(matchedList, r)
	}

	sort.Slice(matchedList, func(i, j int) bool { return matchedList[i].RawPattern() < matchedList[j].RawPattern() })
	return matchedList, nil
}
