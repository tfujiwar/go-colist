package colist

import "fmt"

type ColistEntry struct {
	Pattern string   `json:"pattern"`
	Owners  []string `json:"owners"`
}

// Main opens repository at path, get changed files between the current branch and remote/baseBranch,
// and returns code owners lists that match any of the changed files.
func Main(path string, remote string, baseBranch string) ([]*ColistEntry, error) {
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
