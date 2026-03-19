package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func worktreeRemove(args []string) {
	fs := flag.NewFlagSet("worktree remove", flag.ExitOnError)
	deleteBranch := fs.Bool("delete-branch", false, "also delete the branch associated with the worktree")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip worktree remove [--delete-branch] <submodule> <worktree>")
	}
	fs.Parse(args)

	positional := fs.Args()
	if len(positional) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	submodule := positional[0]
	worktree := positional[1]

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exists, err := submoduleExists(submodule)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !exists {
		fmt.Fprintf(os.Stderr, "submodule %q not found\n", submodule)
		os.Exit(1)
	}

	if err := validateBranchName(worktree); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	root, err := repoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	absWorktreePath := filepath.Join(root, "worktrees", submodule, worktreePathSegment(worktree))
	if _, err := os.Stat(absWorktreePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "worktree %q not found in submodule %q\n", worktree, submodule)
		os.Exit(1)
	}

	submoduleDir := filepath.Join(root, submodule)

	removeCmd := exec.Command("git", "worktree", "remove", absWorktreePath)
	removeCmd.Dir = submoduleDir
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	if err := removeCmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *deleteBranch {
		branchCmd := exec.Command("git", "branch", "-d", worktree)
		branchCmd.Dir = submoduleDir
		branchCmd.Stdout = os.Stdout
		branchCmd.Stderr = os.Stderr
		if err := branchCmd.Run(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
