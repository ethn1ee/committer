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

func GetTrees() (remotes []*git.Remote, workTree *git.Worktree, headTree *object.Tree, err error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	remotes, err = repo.Remotes()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get remotes: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	headTree, err = headCommit.Tree()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get HEAD tree: %w", err)
	}

	workTree, err = repo.Worktree()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	return remotes, workTree, headTree, nil
}

func GetBefore(headTree *object.Tree, path string) (string, error) {
	beforeFile, err := headTree.File(path)
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

func GetAfter(workTree *git.Worktree, path string) (string, error) {
	absPath := filepath.Join(workTree.Filesystem.Root(), path)

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

func Commit(workTree *git.Worktree, msg string) (string, error) {
	hash, err := workTree.Commit(msg, &git.CommitOptions{All: true})
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

func Push(remotes []*git.Remote) error {
	for _, r := range remotes {
		// err := r.Push(nil)
		fmt.Println(r.Config().Name)
	}
	return nil
}
