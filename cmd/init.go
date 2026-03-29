package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Init(args []string) {
	gitDir, err := gitRevParseGitDir()
	if err != nil {
		// Not in any git repo — initialize one.
		if err := runGitInit(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Scaffold .wip.yml if not already present
		if _, err := os.Stat(".wip.yml"); os.IsNotExist(err) {
			cfg := &WipConfig{Refs: map[string]RefConfig{}}
			if err := saveWipConfig(cfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		return
	}

	if strings.TrimSpace(gitDir) == ".git" {
		// Already at the root of a git repo — scaffold .wip.yml if not already present.
		if _, err := os.Stat(".wip.yml"); os.IsNotExist(err) {
			cfg := &WipConfig{Refs: map[string]RefConfig{}}
			if err := saveWipConfig(cfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		return
	}

	// Inside a repo but not at the root.
	fmt.Fprintln(os.Stderr, "not at the root of a git repository")
	os.Exit(1)
}

func gitRevParseGitDir() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stdout = &out
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("not a git repository")
		}
		return "", err
	}
	return out.String(), nil
}

func runGitInit() error {
	cmd := exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("git init failed with exit code %d", exitErr.ExitCode())
		}
		return err
	}
	return nil
}
