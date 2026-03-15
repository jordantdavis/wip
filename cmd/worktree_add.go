package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func worktreeAdd(args []string) {
	fs := flag.NewFlagSet("worktree add", flag.ExitOnError)
	existingBranch := fs.Bool("existing-branch", false, "checkout an existing branch instead of creating a new one")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip worktree add [--existing-branch] <submodule> <worktree>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "args:")
		fmt.Fprintln(os.Stderr, "  submodule    name of the submodule")
		fmt.Fprintln(os.Stderr, "  worktree     name of the worktree (also used as the branch name)")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "flags:")
		fmt.Fprintln(os.Stderr, "  --existing-branch    checkout an existing branch instead of creating a new one")
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

	found, err := submoduleExists(submodule)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !found {
		fmt.Fprintf(os.Stderr, "submodule %q not found\n", submodule)
		os.Exit(1)
	}

	if err := validateWorktreeName(worktree); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	root, err := repoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	absWorktreePath := filepath.Join(root, "worktrees", submodule, worktree)
	submoduleDir := filepath.Join(root, submodule)

	if err := os.MkdirAll(filepath.Join(root, "worktrees", submodule), 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var gitArgs []string
	if *existingBranch {
		gitArgs = []string{"worktree", "add", absWorktreePath, worktree}
	} else {
		gitArgs = []string{"worktree", "add", "-b", worktree, absWorktreePath}
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Dir = submoduleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
