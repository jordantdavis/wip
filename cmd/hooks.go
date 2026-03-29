package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// runHooks executes each hook string via sh -c with WIP_* env vars injected.
// The working directory for each hook is set to worktreePath.
func runHooks(refName, worktreeName, worktreePath, root string, hooks []string) {
	env := append(os.Environ(),
		"WIP_REF_NAME="+refName,
		"WIP_WORKTREE_NAME="+worktreeName,
		"WIP_WORKTREE_PATH="+worktreePath,
		"WIP_ROOT="+root,
	)
	for _, hook := range hooks {
		if hook == "" {
			continue
		}
		hookCmd := exec.Command("sh", "-c", hook)
		hookCmd.Dir = worktreePath
		hookCmd.Stdout = os.Stdout
		hookCmd.Stderr = os.Stderr
		hookCmd.Env = env
		if err := hookCmd.Run(); err != nil {
			fmt.Fprintf(os.Stdout, "✗ %s\n", hook)
		} else {
			fmt.Fprintf(os.Stdout, "✓ %s\n", hook)
		}
	}
}
