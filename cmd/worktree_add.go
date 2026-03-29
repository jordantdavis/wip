package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func worktreeAdd(args []string) {
	fs := flag.NewFlagSet("worktree add", flag.ExitOnError)
	existingBranch := fs.Bool("existing-branch", false, "checkout an existing branch instead of creating a new one")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip worktree add [--existing-branch] <ref> <worktree>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "args:")
		fmt.Fprintln(os.Stderr, "  ref          name of the ref")
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

	cfg, err := requireWipConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	found, err := refExists(submodule)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !found {
		fmt.Fprintf(os.Stderr, "ref %q not found\n", submodule)
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

	hooks := cfg.Refs[submodule].OnWorktreeCreate
	for _, hook := range hooks {
		parts := strings.Fields(hook)
		if len(parts) == 0 {
			continue
		}
		hookCmd := exec.Command(parts[0], parts[1:]...)
		hookCmd.Dir = absWorktreePath
		hookCmd.Stdout = os.Stdout
		hookCmd.Stderr = os.Stderr
		if err := hookCmd.Run(); err != nil {
			fmt.Fprintf(os.Stdout, "✗ %s\n", hook)
		} else {
			fmt.Fprintf(os.Stdout, "✓ %s\n", hook)
		}
	}
}
