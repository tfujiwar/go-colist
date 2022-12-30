package git

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func MkTempRepo() (string, error) {
	dir, err := os.MkdirTemp(os.TempDir(), "go-colist-git-test-")
	if err != nil {
		return "", err
	}

	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		return "", err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	err = os.Mkdir(filepath.Join(dir, ".github"), 0755)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath.Join(dir, ".github", "CODEOWNERS"))
	if err != nil {
		return "", err
	}

	_, err = f.Write([]byte("* all\na.txt @owner-a\nb.txt @owner-b\n"))
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	f, err = os.Create(filepath.Join(dir, "a.txt"))
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	f, err = os.Create(filepath.Join(dir, "b.txt"))
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	_, err = tree.Add(filepath.Join(".github", "CODEOWNERS"))
	if err != nil {
		return "", err
	}

	_, err = tree.Add("a.txt")
	if err != nil {
		return "", err
	}

	_, err = tree.Add("b.txt")
	if err != nil {
		return "", err
	}

	_, err = tree.Commit("add files", &gogit.CommitOptions{})
	if err != nil {
		return "", err
	}

	err = tree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.ReferenceName("feature"),
		Create: true,
	})
	if err != nil {
		return "", err
	}

	err = os.Rename(filepath.Join(dir, "a.txt"), filepath.Join(dir, "c.txt"))
	if err != nil {
		return "", err
	}

	_, err = tree.Add("a.txt")
	if err != nil {
		return "", err
	}

	_, err = tree.Add("c.txt")
	if err != nil {
		return "", err
	}

	_, err = tree.Commit("rename", &gogit.CommitOptions{})
	if err != nil {
		return "", err
	}

	return dir, nil
}

func MkTempRepoLarge() (string, error) {
	numFiles := 100

	dir, err := os.MkdirTemp(os.TempDir(), "go-colist-git-test-")
	if err != nil {
		return "", err
	}

	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		return "", err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	err = os.Mkdir(filepath.Join(dir, ".github"), 0755)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath.Join(dir, ".github", "CODEOWNERS"))
	if err != nil {
		return "", err
	}

	for i := 0; i < numFiles; i++ {
		_, err = f.Write([]byte(fmt.Sprintf("%d/* @owner-%d\n", i, i)))
		if err != nil {
			return "", err
		}
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	_, err = tree.Add(filepath.Join(".github", "CODEOWNERS"))
	if err != nil {
		return "", err
	}

	_, err = tree.Commit("add CODEOWNERS", &gogit.CommitOptions{})
	if err != nil {
		return "", err
	}

	err = tree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.ReferenceName("feature"),
		Create: true,
	})
	if err != nil {
		return "", err
	}

	for i := 0; i < numFiles; i++ {
		err = os.Mkdir(filepath.Join(dir, fmt.Sprintf("%d", i)), 0755)
		if err != nil {
			return "", err
		}

		_, err = os.Create(filepath.Join(dir, fmt.Sprintf("%d", i), "a.txt"))
		if err != nil {
			return "", err
		}

		_, err = tree.Add(filepath.Join(fmt.Sprintf("%d", i), "a.txt"))
		if err != nil {
			return "", err
		}
	}

	_, err = tree.Commit("add files", &gogit.CommitOptions{})
	if err != nil {
		return "", err
	}

	return dir, nil
}

func TestChangedFiles(t *testing.T) {
	dir, err := MkTempRepo()
	defer os.RemoveAll(dir)

	if err != nil {
		t.Error(err)
	}

	files, err := ChangedFiles(dir, "")
	if err != nil {
		t.Error(err)
	}

	if len(files) != 2 {
		t.Errorf("len(files) = %v, want %v", len(files), 2)
	}

	if files[0] != "a.txt" {
		t.Errorf("files[0] = %v, want %v", files[0], "a.txt")
	}

	if files[1] != "c.txt" {
		t.Errorf("files[1] = %v, want %v", files[1], "c.txt")
	}
}
