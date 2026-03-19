package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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
		branch    string
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
			absWorktreePath := filepath.Join(root, "worktrees", submoduleName, worktreeEntry.Name())
			branch := ""
			cmd := exec.Command("git", "-C", absWorktreePath, "branch", "--show-current")
			if out, err := cmd.Output(); err == nil {
				branch = strings.TrimSpace(string(out))
			}
			pairs = append(pairs, pair{submodule: submoduleName, worktree: worktreeEntry.Name(), branch: branch})
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

	subW, wtW := len("SUBMODULE"), len("WORKTREE")
	for _, p := range pairs {
		if len(p.submodule) > subW {
			subW = len(p.submodule)
		}
		if len(p.worktree) > wtW {
			wtW = len(p.worktree)
		}
	}
	fmt.Printf("%-*s  %-*s  %s\n", subW, "SUBMODULE", wtW, "WORKTREE", "BRANCH")
	for _, p := range pairs {
		fmt.Printf("%-*s  %-*s  %s\n", subW, p.submodule, wtW, p.worktree, p.branch)
	}
}
