package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tfujiwar/go-colist/codeowners"
	"github.com/tfujiwar/go-colist/git"
)

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&help, "help", false, "show help")
	flag.Parse()

	if help {
		usage(os.Stdout)
		os.Exit(0)
	}

	var remote string
	var baseBranch string
	switch len(os.Args) {
	case 1:
		remote = ""
		baseBranch = ""
	case 2:
		remote = ""
		baseBranch = os.Args[1]
	case 3:
		remote = os.Args[1]
		baseBranch = os.Args[2]
	default:
		usage(os.Stderr)
		os.Exit(1)
	}

	rules, err := run(".", remote, baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	for _, r := range rules {
		fmt.Printf("%s: %s\n", r.Pattern, strings.Join(r.CodeOwners, ", "))
	}

	os.Exit(0)
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "List GitHub CODEOWNERS of changed files on a current branch\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Usage:\n")
	fmt.Fprintf(w, "  colist                        : compare with remote or local main branch\n")
	fmt.Fprintf(w, "  colist <BASE_BRANCH>          : compare with remote or local <BASE_BRANCH>\n")
	fmt.Fprintf(w, "  colist <REMOTE> <BASE_BRANCH> : compare with <REMOTE>/<BASE_BRANCH>\n")
	fmt.Fprintf(w, "\n")
}

func run(path string, remote string, baseBranch string) ([]*codeowners.Rule, error) {
	files, err := git.ChangedFiles(path, remote, baseBranch)

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
