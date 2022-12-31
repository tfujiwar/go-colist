package main

import (
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
)

// MkTempRepo creates new git repository in dir.
// - Create a repository
// - Commit .github/CODEOWNERS, a.txt, b.txt, and c.txt on master branch
// - Checkout a new branch "feature"
// - Update b.txt and rename c.txt to d.txt
// - Commit b.txt, c.txt, and d.txt
func MkTempRepo(dir string) error {
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		return err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(dir, ".github"), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dir, ".github", "CODEOWNERS"))
	if err != nil {
		return err
	}

	_, err = f.WriteString("* @owner\na.txt @owner-a\nb.txt @owner-b\nc.txt @owner-c\nd.txt @owner-d\n")
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	f, err = os.Create(filepath.Join(dir, "a.txt"))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	f, err = os.Create(filepath.Join(dir, "b.txt"))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	f, err = os.Create(filepath.Join(dir, "c.txt"))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	_, err = tree.Add(filepath.Join(".github", "CODEOWNERS"))
	if err != nil {
		return err
	}

	_, err = tree.Add("a.txt")
	if err != nil {
		return err
	}

	_, err = tree.Add("b.txt")
	if err != nil {
		return err
	}

	_, err = tree.Add("c.txt")
	if err != nil {
		return err
	}

	_, err = tree.Commit("first commit", &gogit.CommitOptions{})
	if err != nil {
		return err
	}

	err = tree.Checkout(&gogit.CheckoutOptions{
		Branch: "feature",
		Create: true,
	})
	if err != nil {
		return err
	}

	f, err = os.OpenFile(filepath.Join(dir, "b.txt"), os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString("foo\n")
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(dir, "c.txt"), filepath.Join(dir, "d.txt"))
	if err != nil {
		return err
	}

	_, err = tree.Add("b.txt")
	if err != nil {
		return err
	}

	_, err = tree.Add("c.txt")
	if err != nil {
		return err
	}

	_, err = tree.Add("d.txt")
	if err != nil {
		return err
	}

	_, err = tree.Commit("second commit", &gogit.CommitOptions{})
	if err != nil {
		return err
	}

	return nil
}

func TestRun(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "go-colist-git-test-")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(dir)

	err = MkTempRepo(dir)
	if err != nil {
		t.Error(err)
		return
	}

	rules, err := run(dir, "", "")
	if err != nil {
		t.Error(err)
		return
	}

	if len(rules) != 3 {
		t.Errorf("len(rules) = %v, want %v", len(rules), 2)
		return
	}
	if len(rules[0].Owners) != 1 {
		t.Errorf("len(rules[0].Owners) = %v, want %v", rules[0].Owners, 1)
		return
	}
	if len(rules[1].Owners) != 1 {
		t.Errorf("len(rules[1].Owners) = %v, want %v", rules[1].Owners, 1)
		return
	}
	if len(rules[2].Owners) != 1 {
		t.Errorf("len(rules[2].Owners) = %v, want %v", rules[2].Owners, 1)
		return
	}

	if rules[0].Pattern != "b.txt" {
		t.Errorf("rules[0].Pattern = %v, want %v", rules[0].Pattern, "b.txt")
	}
	if rules[1].Pattern != "c.txt" {
		t.Errorf("rules[1].Pattern = %v, want %v", rules[1].Pattern, "c.txt")
	}
	if rules[2].Pattern != "d.txt" {
		t.Errorf("rules[2].Pattern = %v, want %v", rules[2].Pattern, "d.txt")
	}

	if rules[0].Owners[0] != "owner-b" {
		t.Errorf("rules[0].Owners[0] = %v, want %v", rules[0].Owners[0], "owner-b")
	}
	if rules[1].Owners[0] != "owner-c" {
		t.Errorf("rules[1].Owners[0] = %v, want %v", rules[1].Owners[0], "owner-c")
	}
	if rules[2].Owners[0] != "owner-d" {
		t.Errorf("rules[2].Owners[0] = %v, want %v", rules[2].Owners[0], "owner-d")
	}
}
