package colist

import "fmt"

func Main(path string, remote string, baseBranch string) ([]*Rule, error) {
	repo, err := NewRepository(path, remote, baseBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to init repo: %w", err)
	}

	cofile, err := repo.OwnersFile()
	if err != nil {
		return nil, fmt.Errorf("failed to open CODEOWNERS: %w", err)
	}

	files, err := repo.ChangedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	rules, err := MatchedRules(cofile, files)
	if err != nil {
		return nil, fmt.Errorf("failed get matched rules: %w", err)
	}

	return rules, nil
}
