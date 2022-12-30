package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var defaultBranches = []string{"main", "master"}

func ChangedFiles(path string, baseBranch string) ([]string, error) {
	opt := gogit.PlainOpenOptions{DetectDotGit: true}
	repo, err := gogit.PlainOpenWithOptions(path, &opt)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	var baseRef *plumbing.Reference
	if baseBranch != "" {
		baseRef, err = repo.Reference(plumbing.ReferenceName("refs/heads/"+baseBranch), false)
		if err != nil {
			return nil, fmt.Errorf("failed to get base branch: %w", err)
		}
	} else {
		for _, b := range defaultBranches {
			baseRef, err = repo.Reference(plumbing.ReferenceName("refs/heads/"+b), false)
			if err == nil {
				break
			}
		}
		if baseRef == nil {
			return nil, fmt.Errorf("failed to get base branch (main or master)")
		}
	}

	baseHead, err := repo.CommitObject(baseRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get base commit: %w", err)
	}

	baseCommits, err := commit.MergeBase(baseHead)
	if err != nil || len(baseCommits) == 0 {
		return nil, fmt.Errorf("failed to get merge base: %w", err)
	}

	baseTree, err := baseCommits[0].Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	changes, err := object.DiffTree(baseTree, tree)
	if err != nil {
		return nil, fmt.Errorf("failed to get diffs: %w", err)
	}

	fileset := make(map[string]struct{})
	for _, c := range changes {
		if c.From.Name != "" {
			fileset[c.From.Name] = struct{}{}
		}
		if c.To.Name != "" {
			fileset[c.To.Name] = struct{}{}
		}
	}

	files := make([]string, 0)
	for f := range fileset {
		files = append(files, f)
	}

	return files, nil
}
