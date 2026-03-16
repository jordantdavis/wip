package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func worktreeList(_ []string) {
	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	root, err := repoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	worktreesDir := filepath.Join(root, "worktrees")
	if _, err := os.Stat(worktreesDir); os.IsNotExist(err) {
		fmt.Println("no worktrees found")
		os.Exit(0)
	}

	submoduleDirs, err := os.ReadDir(worktreesDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	type pair struct {
		submodule string
		worktree  string
	}

	var pairs []pair

	for _, submoduleEntry := range submoduleDirs {
		if !submoduleEntry.IsDir() {
			continue
		}
		submoduleName := submoduleEntry.Name()
		worktreeDirs, err := os.ReadDir(filepath.Join(worktreesDir, submoduleName))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, worktreeEntry := range worktreeDirs {
			if !worktreeEntry.IsDir() {
				continue
			}
			pairs = append(pairs, pair{submodule: submoduleName, worktree: worktreeEntry.Name()})
		}
	}

	if len(pairs) == 0 {
		fmt.Println("no worktrees found")
		os.Exit(0)
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].submodule != pairs[j].submodule {
			return pairs[i].submodule < pairs[j].submodule
		}
		return pairs[i].worktree < pairs[j].worktree
	})

	for _, p := range pairs {
		fmt.Printf("%s  %s\n", p.submodule, p.worktree)
	}
}
