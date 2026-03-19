package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Worktree(args []string) {
	if len(args) < 1 {
		worktreeUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		worktreeAdd(args[1:])
	case "list":
		worktreeList(args[1:])
	case "remove":
		worktreeRemove(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown worktree command: %s\n", args[0])
		worktreeUsage()
		os.Exit(1)
	}
}

func worktreeUsage() {
	fmt.Fprintln(os.Stderr, "usage: wip worktree <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  add <submodule> <worktree>      create a new worktree in a submodule")
	fmt.Fprintln(os.Stderr, "  list                            list all worktrees")
	fmt.Fprintln(os.Stderr, "  remove <submodule> <worktree>   remove a worktree from a submodule")
}

func validateBranchName(name string) error {
	cmd := exec.Command("git", "check-ref-format", "--branch", name)
	cmd.Dir, _ = os.Getwd()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("invalid branch name %q", name)
	}
	return nil
}

func worktreePathSegment(branchName string) string {
	return strings.ReplaceAll(branchName, "/", "-")
}

func repoRoot() (string, error) {
	return os.Getwd()
}
