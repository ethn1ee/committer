package committer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ethn1ee/committer/internal/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/utils/binary"
	"google.golang.org/genai"
)

type Diff struct {
	Path   string
	Before string
	After  string
}

const binaryFile = "binary file"

func Generate(cfg *config.Config) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	chat, err := client.Chats.Create(ctx, "gemini-2.0-flash", nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini chat: %w", err)
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	headTree, err := headCommit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD tree: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	_, err = worktree.Add(".")
	if err != nil {
		return "", fmt.Errorf("failed to add changes to staging: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree status: %w", err)
	}

	if status.IsClean() {
		return "No changes to commit", nil
	}

	diffs := make([]*Diff, 0, len(status))

	for path, fs := range status {
		switch fs.Staging {
		case git.Modified:
			before, err := getBefore(headTree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get before contents for file %s: %w", path, err)
			}

			after, err := getAfter(worktree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get after contents for file %s: %w", path, err)
			}

			diffs = append(diffs, &Diff{
				Path:   path,
				Before: before,
				After:  after,
			})

		case git.Added:
			after, err := getAfter(worktree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get after contents for file %s: %w", path, err)
			}

			diffs = append(diffs, &Diff{
				Path:   path,
				Before: "",
				After:  after,
			})

		case git.Deleted:
			before, err := getBefore(headTree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get before contents for file %s: %w", path, err)
			}

			diffs = append(diffs, &Diff{
				Path:   path,
				Before: before,
				After:  "",
			})

		case git.Renamed:
			// fmt.Fprintf(os.Stdout, "Renamed: %s\n", path)
		case git.Untracked:
			// fmt.Fprintf(os.Stdout, "Untracked: %s\n", path)
		default:
			// fmt.Fprintf(os.Stdout, "Other change (%s): %s\n", string(fs.Worktree), path)
		}
	}

	diffsByte, err := json.Marshal(diffs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal diffs: %w", err)
	}

	prompt := fmt.Sprintf(
		`Create a concise git commit message based on the following:
status:%s
beforeAfter:%s
rules:
- The message contains two parts: header and body
- The header (first line of the commit) should be in one of the following formats:
	- feat: <short description>
	- enhancement: <short description>
	- fix: <short description>
	- docs: <short description>
- Do not capitalize the first letter of the description.
- Do not write multi-line messages
- The body (optional, after the header) should provide additional context and details about the changes made in a bullet list.
`,
		status.String(),
		string(diffsByte),
	)

	res, err := chat.Send(ctx, &genai.Part{
		Text: prompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message to Gemini: %w", err)
	}

	return res.Text(), nil
}

func getBefore(headtree *object.Tree, path string) (string, error) {
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

func getAfter(worktree *git.Worktree, path string) (string, error) {
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

	var afterBuffer bytes.Buffer
	_, err = io.Copy(&afterBuffer, afterFile)
	if err != nil {
		return "", fmt.Errorf("failed to read contents of file %s: %w", absPath, err)
	}

	after := afterBuffer.String()

	return after, nil
}
