package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hmarr/codeowners"
	"github.com/tfujiwar/go-colist/git"
)

func main() {
	var baseBranch string
	switch len(os.Args) {
	case 1:
	case 2:
		baseBranch = os.Args[1]
	default:
		fmt.Fprintf(os.Stderr, "usage: %s <base-branch>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	if err := run(baseBranch); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(baseBranch string) error {
	files, err := git.ChangedFiles(".", baseBranch)

	if err != nil {
		return fmt.Errorf("failed to get changed files: %w", err)
	}

	cofile, err := os.Open(filepath.Join(".", ".github/CODEOWNERS"))
	if err != nil {
		return fmt.Errorf("failed to open CODEOWNERS: %w", err)
	}

	ruleset, err := codeowners.ParseFile(cofile)
	if err != nil {
		return fmt.Errorf("failed to parse CODEOWNERS: %w", err)
	}

	matched := make(map[string]*codeowners.Rule)
	for _, f := range files {
		rule, err := ruleset.Match(f)
		if err != nil {
			return fmt.Errorf("failed to match CODEOWNERS rule: %w", err)
		}
		matched[rule.RawPattern()] = rule
	}

	for _, r := range matched {
		fmt.Printf("%v: %v\n", r.RawPattern(), r.Owners)
	}

	return nil
}
