package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupLaunchEnv creates a temp dir with a .git entry, a .wip.yml, and the
// worktree directory, changes the working directory to root, and returns a
// restore function that must be deferred by the caller.
func setupLaunchEnv(t *testing.T, submodule, worktree string, cfg *WipConfig) (root string, restore func()) {
	t.Helper()
	root = t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".git"), 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	worktreeDir := filepath.Join(root, "worktrees", submodule, worktreePathSegment(worktree))
	if err := os.MkdirAll(worktreeDir, 0755); err != nil {
		t.Fatalf("mkdir worktree: %v", err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := saveWipConfig(cfg); err != nil {
		t.Fatalf("saveWipConfig: %v", err)
	}

	return root, func() { os.Chdir(orig) }
}

// captureStdout replaces os.Stdout with a pipe and returns the captured string
// after calling f. It restores os.Stdout before returning.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	tmp := make([]byte, 4096)
	for {
		n, readErr := r.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
		}
		if readErr != nil {
			break
		}
	}
	r.Close()
	return buf.String()
}

// TestWorktreeLaunch_HooksRunSuccessfully verifies that a hook that exits 0
// produces a ✓ line in stdout.
func TestWorktreeLaunch_HooksRunSuccessfully(t *testing.T) {
	cfg := &WipConfig{
		Refs: map[string]RefConfig{
			"myservice": {OnWorktreeLaunch: []string{"echo hello"}},
		},
	}
	_, restore := setupLaunchEnv(t, "myservice", "my-feature", cfg)
	defer restore()

	output := captureStdout(t, func() {
		worktreeLaunch([]string{"myservice", "my-feature"})
	})

	if !strings.Contains(output, "✓ echo hello") {
		t.Errorf("expected '✓ echo hello' in output, got: %q", output)
	}
}

// TestWorktreeLaunch_HookFailurePrintsX verifies that a hook that exits non-zero
// produces a ✗ line and execution continues to subsequent hooks.
func TestWorktreeLaunch_HookFailurePrintsX(t *testing.T) {
	cfg := &WipConfig{
		Refs: map[string]RefConfig{
			"myservice": {OnWorktreeLaunch: []string{"false", "echo after"}},
		},
	}
	_, restore := setupLaunchEnv(t, "myservice", "my-feature", cfg)
	defer restore()

	output := captureStdout(t, func() {
		worktreeLaunch([]string{"myservice", "my-feature"})
	})

	if !strings.Contains(output, "✗ false") {
		t.Errorf("expected '✗ false' in output, got: %q", output)
	}
	if !strings.Contains(output, "✓ echo after") {
		t.Errorf("expected subsequent hook '✓ echo after' to still run, got: %q", output)
	}
}

// TestWorktreeLaunch_MultipleHooksRunInOrder verifies that multiple hooks run
// in the configured order.
func TestWorktreeLaunch_MultipleHooksRunInOrder(t *testing.T) {
	cfg := &WipConfig{
		Refs: map[string]RefConfig{
			"myservice": {OnWorktreeLaunch: []string{"echo first", "echo second"}},
		},
	}
	_, restore := setupLaunchEnv(t, "myservice", "my-feature", cfg)
	defer restore()

	output := captureStdout(t, func() {
		worktreeLaunch([]string{"myservice", "my-feature"})
	})

	firstPos := strings.Index(output, "✓ echo first")
	secondPos := strings.Index(output, "✓ echo second")
	if firstPos == -1 || secondPos == -1 {
		t.Fatalf("expected both hooks in output, got: %q", output)
	}
	if firstPos > secondPos {
		t.Errorf("expected 'echo first' before 'echo second' in output, got: %q", output)
	}
}

// TestWorktreeLaunch_HooksCwdIsWorktreeDir verifies that hooks run with their
// working directory set to the worktree directory by writing a sentinel file
// from within the hook and checking it appears in the right location.
func TestWorktreeLaunch_HooksCwdIsWorktreeDir(t *testing.T) {
	cfg := &WipConfig{
		Refs: map[string]RefConfig{
			"myservice": {OnWorktreeLaunch: []string{"touch sentinel.txt"}},
		},
	}
	root, restore := setupLaunchEnv(t, "myservice", "my-feature", cfg)
	defer restore()

	captureStdout(t, func() {
		worktreeLaunch([]string{"myservice", "my-feature"})
	})

	sentinel := filepath.Join(root, "worktrees", "myservice", "my-feature", "sentinel.txt")
	if _, err := os.Stat(sentinel); os.IsNotExist(err) {
		t.Errorf("expected sentinel.txt to be created in worktree dir, but it was not found at %s", sentinel)
	}
}
