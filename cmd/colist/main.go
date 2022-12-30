package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tfujiwar/go-colist/codeowners"
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

	rules, err := run(".", baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	for _, r := range rules {
		fmt.Printf("%s: %s\n", r.Pattern, strings.Join(r.CodeOwners, ", "))
	}

	os.Exit(0)
}

func run(path string, baseBranch string) ([]*codeowners.Rule, error) {
	files, err := git.ChangedFiles(path, baseBranch)

	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	cofile, err := os.Open(filepath.Join(path, ".github/CODEOWNERS"))
	if err != nil {
		return nil, fmt.Errorf("failed to open CODEOWNERS: %w", err)
	}

	rules, err := codeowners.MatchedRules(cofile, files)
	if err != nil {
		return nil, fmt.Errorf("failed get matched rules: %w", err)
	}

	return rules, nil
}
