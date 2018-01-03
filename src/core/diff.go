package util

import (
	"bytes"
	"os/exec"
	"strings"
)

// Diff get diff files between two git commits
func Diff(commit1, commit2 string) ([]string, error) {
	gitDiff := exec.Command("git", "diff", "--name-only", commit1, commit2)
	var out bytes.Buffer
	gitDiff.Stdout = &out
	if err := gitDiff.Start(); err != nil {
		return nil, err
	}

	if err := gitDiff.Wait(); err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(out.String()), "\n"), nil
}
