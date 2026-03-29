package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func refRemove(args []string) {
	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: wip ref remove <name>")
		os.Exit(1)
	}

	name := args[0]

	exists, err := refExists(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !exists {
		fmt.Fprintf(os.Stderr, "ref %q not found\n", name)
		os.Exit(1)
	}

	// Step 1: git submodule deinit -f <name>
	deinit := exec.Command("git", "submodule", "deinit", "-f", name)
	deinit.Stdout = os.Stdout
	deinit.Stderr = os.Stderr
	if err := deinit.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "deinit failed")
		os.Exit(1)
	}

	// Step 2: git rm -f <name>
	gitRm := exec.Command("git", "rm", "-f", name)
	gitRm.Stdout = os.Stdout
	gitRm.Stderr = os.Stderr
	if err := gitRm.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "git rm failed")
		os.Exit(1)
	}

	// Step 3: remove .git/modules/<name>
	if err := os.RemoveAll(filepath.Join(".git", "modules", name)); err != nil {
		fmt.Fprintf(os.Stderr, "failed to remove .git/modules/%s: %v\n", name, err)
		os.Exit(1)
	}

	// Remove from .wip.yml if present
	cfg, err := requireWipConfig()
	if err != nil {
		// .wip.yml not found — skip silently
		return
	}
	if cfg.Refs != nil {
		if _, ok := cfg.Refs[name]; ok {
			delete(cfg.Refs, name)
			if err := saveWipConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "failed to update .wip.yml: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
