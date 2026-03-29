package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func worktreeLaunch(args []string) {
	fs := flag.NewFlagSet("worktree launch", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip worktree launch <ref> <worktree>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "args:")
		fmt.Fprintln(os.Stderr, "  ref          name of the ref")
		fmt.Fprintln(os.Stderr, "  worktree     name of the worktree")
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

	root, err := repoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	absWorktreePath := filepath.Join(root, "worktrees", submodule, worktreePathSegment(worktree))
	if _, err := os.Stat(absWorktreePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "worktree %q not found for ref %q\n", worktree, submodule)
		os.Exit(1)
	}

	hooks := cfg.Refs[submodule].OnWorktreeLaunch
	if len(hooks) == 0 {
		fmt.Fprintf(os.Stdout, "no on-worktree-launch hooks configured for %s\n", submodule)
		os.Exit(0)
	}

	runHooks(submodule, worktree, absWorktreePath, root, hooks)
}
