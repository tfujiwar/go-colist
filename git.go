package colist

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// newRepository opens a git repository at the path.
func newRepository(path string) (*gogit.Repository, error) {
	opt := gogit.PlainOpenOptions{DetectDotGit: true}
	repo, err := gogit.PlainOpenWithOptions(path, &opt)
	if err != nil {
		return nil, fmt.Errorf("open repo at %s: %w", path, err)
	}

	return repo, nil
}

// currentCommitAndTree returns a commit object and a tree object of the repo.
func currentCommitAndTree(repo *gogit.Repository) (*object.Commit, *object.Tree, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, nil, fmt.Errorf("HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("commit at %s: %w", ref.Hash(), err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, nil, fmt.Errorf("tree at %s: %w", ref.Hash(), err)
	}

	return commit, tree, nil
}

// baseCommitAndTree returns a commit object and a tree object at the remote and the branch of the repo.
func baseCommitAndTree(repo *gogit.Repository, remote, branch string) (*object.Commit, *object.Tree, error) {
	refs := baseRefCandidates(remote, branch)

	var ref *plumbing.Reference
	for _, r := range refs {
		var err error
		ref, err = repo.Reference(plumbing.ReferenceName(r), false)
		if err == nil {
			log.Printf("[DEBUG] found baseRef : %s\n", r)
			break
		}
		log.Printf("[DEBUG] tried baseRef : %s\n", r)
	}

	if ref == nil {
		return nil, nil, fmt.Errorf("cannot find any of refs: %s", strings.Join(refs, ", "))
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("commit at %s: %w", ref.Hash(), err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, nil, fmt.Errorf("tree at %s: %w", ref.Hash(), err)
	}

	return commit, tree, nil
}

func baseRefCandidates(remote, branch string) []string {
	if remote == "" {
		if branch == "" {
			return []string{
				"refs/remote/origin/main",
				"refs/remote/origin/master",
				"refs/heads/main",
				"refs/heads/master",
			}
		} else {
			return []string{
				"refs/remote/origin/" + branch,
				"refs/heads/" + branch,
			}
		}
	} else {
		if branch == "" {
			return []string{
				"refs/remote/" + remote + "/main",
				"refs/remote/" + remote + "/master",
			}
		} else {
			return []string{
				"refs/remote/" + remote + "/" + branch,
			}
		}
	}
}

// codeOwnersFile returns a reader for CODEOWNERS file.
func codeOwnersFile(tree *object.Tree) (io.Reader, error) {
	path := ".github/CODEOWNERS"
	f, err := tree.File(path)
	if err != nil {
		return nil, fmt.Errorf("at %s: %w", path, err)
	}

	return f.Reader()
}

// mergeBaseCommitAndTree returns a commit object and a tree object of the given commits c1 and c2.
func mergeBaseCommitAndTree(c1, c2 *object.Commit) (*object.Commit, *object.Tree, error) {
	commits, err := c1.MergeBase(c2)
	if err != nil || len(commits) == 0 {
		return nil, nil, fmt.Errorf("merge base of %s and %s : %w", c1.Hash, c2.Hash, err)
	}

	tree, err := commits[0].Tree()
	if err != nil {
		return nil, nil, fmt.Errorf("tree at %s: %w", commits[0].Hash, err)
	}

	return commits[0], tree, nil
}

// changedFiles returns a list of files changed between the tree object from and to
func changedFiles(to, from *object.Tree) ([]string, error) {
	changes, err := object.DiffTree(to, from)
	if err != nil {
		return nil, fmt.Errorf("diff: %w", err)
	}

	files := make([]string, 0)
	for _, c := range changes {
		if c.From.Name == c.To.Name {
			log.Printf("[DEBUG] updated file: %s\n", c.From.Name)
			files = append(files, c.From.Name)
		} else {
			if c.From.Name != "" {
				log.Printf("[DEBUG] created file: %s\n", c.From.Name)
				files = append(files, c.From.Name)
			}
			if c.To.Name != "" {
				log.Printf("[DEBUG] deleted file: %s\n", c.To.Name)
				files = append(files, c.To.Name)
			}
		}
	}

	sort.Slice(files, func(i, j int) bool { return files[i] < files[j] })
	return files, nil
}
