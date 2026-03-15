package cmd

import (
	"fmt"
	"os"
	"regexp"
)

var worktreeNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

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

func validateWorktreeName(name string) error {
	if !worktreeNameRe.MatchString(name) {
		return fmt.Errorf("worktree name %q must match [a-zA-Z0-9_-]+", name)
	}
	return nil
}

func repoRoot() (string, error) {
	return os.Getwd()
}
