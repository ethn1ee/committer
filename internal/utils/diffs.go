package utils

import (
	"fmt"

	"github.com/ethn1ee/committer/internal/models"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func GetDiffs(headTree *object.Tree, workTree *git.Worktree) ([]*models.Diff, error) {
	status, err := GetStatus(workTree)
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	diffs := make([]*models.Diff, 0, len(status))

	for path, fs := range status {
		d := &models.Diff{
			Path:       path,
			StatusCode: string(fs.Staging),
		}

		switch fs.Staging {
		case git.Modified:
			before, err := GetBefore(headTree, path)
			if err != nil {
				return nil, fmt.Errorf("failed to get before contents for file %s: %w", path, err)
			}

			after, err := GetAfter(workTree, path)
			if err != nil {
				return nil, fmt.Errorf("failed to get after contents for file %s: %w", path, err)
			}

			d.Before = before
			d.After = after

		case git.Added:
			after, err := GetAfter(workTree, path)
			if err != nil {
				return nil, fmt.Errorf("failed to get after contents for file %s: %w", path, err)
			}

			d.Before = ""
			d.After = after

		case git.Deleted:
			before, err := GetBefore(headTree, path)
			if err != nil {
				return nil, fmt.Errorf("failed to get before contents for file %s: %w", path, err)
			}

			d.Before = before
			d.After = ""
		}

		diffs = append(diffs, d)
	}

	return diffs, nil
}
