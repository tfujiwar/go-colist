package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tfujiwar/go-colist/codeowners"
	"github.com/tfujiwar/go-colist/format"
	"github.com/tfujiwar/go-colist/git"
)

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&help, "help", false, "show help")

	var dir string
	flag.StringVar(&dir, "d", ".", "path to repository directory")
	flag.StringVar(&dir, "dir", ".", "path to repository directory")

	flag.Parse()

	if help {
		usage(os.Stdout)
		os.Exit(0)
	}

	args := flag.Args()

	var remote string
	var baseBranch string
	switch len(args) {
	case 0:
		remote = ""
		baseBranch = ""
	case 1:
		remote = ""
		baseBranch = args[1]
	case 2:
		remote = args[1]
		baseBranch = args[2]
	default:
		usage(os.Stderr)
		os.Exit(1)
	}

	rules, err := run(dir, remote, baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	format.TextWithIndent(rules, os.Stdout)
	os.Exit(0)
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "List GitHub CODEOWNERS of changed files on a current branch\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Usage:\n")
	fmt.Fprintf(w, "  colist [flags]                        : compare with remote or local main branch\n")
	fmt.Fprintf(w, "  colist [flags] <BASE_BRANCH>          : compare with remote or local <BASE_BRANCH>\n")
	fmt.Fprintf(w, "  colist [flags] <REMOTE> <BASE_BRANCH> : compare with <REMOTE>/<BASE_BRANCH>\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Flags:\n")
	fmt.Fprintf(w, "  -d, --dir  : path to repository directory\n")
	fmt.Fprintf(w, "  -h, --help : show this message\n")
	fmt.Fprintf(w, "\n")
}

func run(path string, remote string, baseBranch string) ([]*codeowners.Rule, error) {
	repo, err := git.NewRepository(path, remote, baseBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to init repo: %w", err)
	}

	cofile, err := repo.CodeOwnersFile()
	if err != nil {
		return nil, fmt.Errorf("failed to open CODEOWNERS: %w", err)
	}

	files, err := repo.ChangedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	rules, err := codeowners.MatchedRules(cofile, files)
	if err != nil {
		return nil, fmt.Errorf("failed get matched rules: %w", err)
	}

	return rules, nil
}
