package colist

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type ColistEntry struct {
	Pattern string   `json:"pattern"`
	Owners  []string `json:"owners"`
}

// Run runs a main logic of colist
func Run() error {
	var output string
	flag.StringVar(&output, "o", "text", "")
	flag.StringVar(&output, "output", "text", "")

	var dir string
	flag.StringVar(&dir, "d", ".", "")
	flag.StringVar(&dir, "dir", ".", "")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&verbose, "verbose", false, "")

	var help bool
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")

	flag.Parse()

	if help {
		fmt.Printf("A command to show GitHub code owners of files that were modified in a current branch\n")
		fmt.Printf("\n")
		fmt.Printf("Usage:\n")
		fmt.Printf("  colist [flags]                        : compare with remote or local main branch\n")
		fmt.Printf("  colist [flags] <BASE_BRANCH>          : compare with remote or local <BASE_BRANCH>\n")
		fmt.Printf("  colist [flags] <REMOTE> <BASE_BRANCH> : compare with <REMOTE>/<BASE_BRANCH>\n")
		fmt.Printf("\n")
		fmt.Printf("Flags:\n")
		fmt.Printf("  -o, --output text|json : output format\n")
		fmt.Printf("  -d, --dir <DIR>        : repository directory\n")
		fmt.Printf("  -v, --verbose          : show debug log\n")
		fmt.Printf("  -h, --help             : show this message\n")
		fmt.Printf("\n")
		return nil
	}

	// Use log package for debug log
	if verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}

	var formatFunc func([]*ColistEntry, io.Writer) error
	switch output {
	case "text":
		formatFunc = outputText
	case "json":
		formatFunc = outputJson
	default:
		return fmt.Errorf("[ERROR] not supported output: select \"text\" or \"json\"")
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
		baseBranch = args[0]
	case 2:
		remote = args[0]
		baseBranch = args[1]
	default:
		return fmt.Errorf("[ERROR] too many args")
	}

	log.Printf("[DEBUG] output     : %v\n", output)
	log.Printf("[DEBUG] dir        : %v\n", dir)
	log.Printf("[DEBUG] verbose    : %v\n", verbose)
	log.Printf("[DEBUG] help       : %v\n", help)
	log.Printf("[DEBUG] remote     : %v\n", remote)
	log.Printf("[DEBUG] baseBranch : %v\n", baseBranch)

	colists, err := run(dir, remote, baseBranch)
	if err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	err = formatFunc(colists, os.Stdout)
	if err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}

// run opens repository at path, get changed files between the current branch and remote/baseBranch,
// and returns code owners lists that match any of the changed files.
func run(path string, remote string, baseBranch string) ([]*ColistEntry, error) {
	repo, err := newRepository(path)
	if err != nil {
		return nil, fmt.Errorf("failed to init repo: %w", err)
	}

	currentCommit, currentTree, err := currentCommitAndTree(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get current commit and tree: %w", err)
	}

	baseCommit, baseTree, err := baseCommitAndTree(repo, remote, baseBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get base commit and tree: %w", err)
	}

	_, mbTree, err := mergeBaseCommitAndTree(currentCommit, baseCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge base commit and tree: %w", err)
	}

	files, err := changedFiles(currentTree, mbTree)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	cofile, err := codeOwnersFile(baseTree)
	if err != nil {
		return nil, fmt.Errorf("failed to open CODEOWNERS: %w", err)
	}

	colists, err := codeOwnersLists(cofile, files)
	if err != nil {
		return nil, fmt.Errorf("failed get code owners lists: %w", err)
	}

	return colists, nil
}
