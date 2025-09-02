package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/utils/binary"
)

const binaryFile = "binary file"

func GetTrees() (worktree *git.Worktree, headtree *object.Tree, err error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	headtree, err = headCommit.Tree()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get HEAD tree: %w", err)
	}

	worktree, err = repo.Worktree()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	return worktree, headtree, nil
}

func GetBefore(headtree *object.Tree, path string) (string, error) {
	beforeFile, err := headtree.File(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file %s from HEAD tree: %w", path, err)
	}

	isBin, err := beforeFile.IsBinary()
	if err != nil {
		return "", fmt.Errorf("failed to check if file %s is binary: %w", path, err)
	}
	if isBin {
		return binaryFile, nil
	}

	before, err := beforeFile.Contents()
	if err != nil {
		return "", fmt.Errorf("failed to get contents of file %s: %w", path, err)
	}

	return before, nil
}

func GetAfter(worktree *git.Worktree, path string) (string, error) {
	absPath := filepath.Join(worktree.Filesystem.Root(), path)

	afterFile, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", absPath, err)
	}
	defer afterFile.Close()

	isBin, err := binary.IsBinary(afterFile)
	if isBin {
		return binaryFile, nil
	}

	_, err = afterFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek to beginning of file %s: %w", absPath, err)
	}

	var afterBuffer bytes.Buffer
	_, err = io.Copy(&afterBuffer, afterFile)
	if err != nil {
		return "", fmt.Errorf("failed to read contents of file %s: %w", absPath, err)
	}

	after := afterBuffer.String()

	return after, nil
}
