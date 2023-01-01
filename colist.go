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
		usage(os.Stdout)
		os.Exit(0)
	}

	if verbose {
		log.SetFlags(0)
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}

	var formatFunc func([]*ColistEntry, io.Writer) error
	switch output {
	case "text":
		formatFunc = OutputText
	case "json":
		formatFunc = OutputJson
	default:
		return fmt.Errorf("not supported output: select \"text\" or \"json\"")
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
		usage(os.Stderr)
		os.Exit(1)
	}

	log.Printf("[DEBUG] output     : %v\n", output)
	log.Printf("[DEBUG] dir        : %v\n", dir)
	log.Printf("[DEBUG] verbose    : %v\n", verbose)
	log.Printf("[DEBUG] help       : %v\n", help)
	log.Printf("[DEBUG] remote     : %v\n", remote)
	log.Printf("[DEBUG] baseBranch : %v\n", baseBranch)

	colists, err := run(dir, remote, baseBranch)
	if err != nil {
		return err
	}

	err = formatFunc(colists, os.Stdout)
	if err != nil {
		return err
	}

	return nil
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
	fmt.Fprintf(w, "  -o, --output text|json : output format\n")
	fmt.Fprintf(w, "  -d, --dir <DIR>        : repository directory\n")
	fmt.Fprintf(w, "  -v, --verbose          : show debug log\n")
	fmt.Fprintf(w, "  -h, --help             : show this message\n")
	fmt.Fprintf(w, "\n")
}

// run opens repository at path, get changed files between the current branch and remote/baseBranch,
// and returns code owners lists that match any of the changed files.
func run(path string, remote string, baseBranch string) ([]*ColistEntry, error) {
	repo, err := NewRepository(path)
	if err != nil {
		return nil, fmt.Errorf("failed to init repo: %w", err)
	}

	currentCommit, currentTree, err := CurrentCommitAndTree(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get current commit and tree: %w", err)
	}

	baseCommit, baseTree, err := BaseCommitAndTree(repo, remote, baseBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get base commit and tree: %w", err)
	}

	_, mbTree, err := MergeBaseCommitAndTree(currentCommit, baseCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge base commit and tree: %w", err)
	}

	files, err := ChangedFiles(currentTree, mbTree)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	cofile, err := CodeOwnersFile(baseTree)
	if err != nil {
		return nil, fmt.Errorf("failed to open CODEOWNERS: %w", err)
	}

	colists, err := CodeOwnersLists(cofile, files)
	if err != nil {
		return nil, fmt.Errorf("failed get code owners lists: %w", err)
	}

	return colists, nil
}
