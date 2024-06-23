package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func FileList(repository, revision string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", revision, "--name-only")
	cmd.Dir = repository

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-tree -r %s: %v", revision, err)
	}

	files, err := splitLines(string(out))
	if err != nil {
		return nil, fmt.Errorf("split lines: %v", err)
	}

	return files, nil
}

func splitLines(s string) ([]string, error) {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
