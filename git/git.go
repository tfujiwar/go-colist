package git

import (
	"fmt"
	"io"
	"sort"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repository struct {
	currentTree *object.Tree
	latestTree  *object.Tree
	baseTree    *object.Tree
}

func NewRepository(path string, remote string, baseBranch string) (*Repository, error) {
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
		return nil, fmt.Errorf("failed to get current tree: %w", err)
	}

	var refs []string
	if remote == "" {
		if baseBranch == "" {
			refs = []string{
				"refs/remote/origin/main",
				"refs/remote/origin/master",
				"refs/heads/main",
				"refs/heads/master",
			}
		} else {
			refs = []string{
				"refs/remote/origin/" + baseBranch,
				"refs/heads/" + baseBranch,
			}
		}
	} else {
		if baseBranch == "" {
			refs = []string{
				"refs/remote/" + remote + "/main",
				"refs/remote/" + remote + "/master",
			}
		} else {
			refs = []string{
				"refs/remote/" + remote + "/" + baseBranch,
			}
		}
	}

	var baseRef *plumbing.Reference
	for _, r := range refs {
		baseRef, err = repo.Reference(plumbing.ReferenceName(r), false)
		if err == nil {
			break
		}
	}
	if baseRef == nil {
		return nil, fmt.Errorf("failed to get any of refs: %s", strings.Join(refs, ", "))
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
		return nil, fmt.Errorf("failed to get base tree: %w", err)
	}

	latestTree, err := baseHead.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest tree: %w", err)
	}

	return &Repository{
		currentTree: tree,
		latestTree:  latestTree,
		baseTree:    baseTree,
	}, nil
}

func (r *Repository) CodeOwnersFile() (io.Reader, error) {
	f, err := r.currentTree.File(".github/CODEOWNERS")
	if err != nil {
		return nil, fmt.Errorf("failed to find CODEOWNERS: %w", err)
	}

	reader, err := f.Reader()
	if err != nil {
		return nil, fmt.Errorf("failed to read CODEOWNERS: %w", err)
	}

	return reader, nil
}

func (r *Repository) ChangedFiles() ([]string, error) {
	changes, err := object.DiffTree(r.baseTree, r.currentTree)
	if err != nil {
		return nil, fmt.Errorf("failed to get diffs: %w", err)
	}

	files := make([]string, 0)
	for _, c := range changes {
		if c.From.Name == c.To.Name {
			files = append(files, c.From.Name)
		} else {
			if c.From.Name != "" {
				files = append(files, c.From.Name)
			}
			if c.To.Name != "" {
				files = append(files, c.To.Name)
			}
		}
	}

	sort.Slice(files, func(i, j int) bool { return files[i] < files[j] })
	return files, nil
}
