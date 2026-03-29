package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunHooks_WipRefName verifies WIP_REF_NAME is available in hooks.
func TestRunHooks_WipRefName(t *testing.T) {
	dir := t.TempDir()
	output := captureStdout(t, func() {
		runHooks("myref", "myworktree", dir, "/fake/root", []string{"echo $WIP_REF_NAME"})
	})
	if !strings.Contains(output, "myref") {
		t.Errorf("expected WIP_REF_NAME=myref in hook output, got: %q", output)
	}
}

// TestRunHooks_WipWorktreeName verifies WIP_WORKTREE_NAME is available in hooks.
func TestRunHooks_WipWorktreeName(t *testing.T) {
	dir := t.TempDir()
	output := captureStdout(t, func() {
		runHooks("myref", "myworktree", dir, "/fake/root", []string{"echo $WIP_WORKTREE_NAME"})
	})
	if !strings.Contains(output, "myworktree") {
		t.Errorf("expected WIP_WORKTREE_NAME=myworktree in hook output, got: %q", output)
	}
}

// TestRunHooks_WipWorktreePath verifies WIP_WORKTREE_PATH is the absolute worktree path.
func TestRunHooks_WipWorktreePath(t *testing.T) {
	dir := t.TempDir()
	absDir, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	output := captureStdout(t, func() {
		runHooks("myref", "myworktree", absDir, "/fake/root", []string{"echo $WIP_WORKTREE_PATH"})
	})
	if !strings.Contains(output, absDir) {
		t.Errorf("expected WIP_WORKTREE_PATH=%q in hook output, got: %q", absDir, output)
	}
}

// TestRunHooks_WipRoot verifies WIP_ROOT is the repo root.
func TestRunHooks_WipRoot(t *testing.T) {
	dir := t.TempDir()
	root := "/fake/root/dir"
	output := captureStdout(t, func() {
		runHooks("myref", "myworktree", dir, root, []string{"echo $WIP_ROOT"})
	})
	if !strings.Contains(output, root) {
		t.Errorf("expected WIP_ROOT=%q in hook output, got: %q", root, output)
	}
}

// TestRunHooks_CompoundCommand verifies compound commands via && execute both parts.
func TestRunHooks_CompoundCommand(t *testing.T) {
	dir := t.TempDir()
	output := captureStdout(t, func() {
		runHooks("myref", "myworktree", dir, "/fake/root", []string{"echo a && echo b"})
	})
	if !strings.Contains(output, "a") {
		t.Errorf("expected 'a' in output, got: %q", output)
	}
	if !strings.Contains(output, "b") {
		t.Errorf("expected 'b' in output, got: %q", output)
	}
}

// TestRunHooks_WritesFileInWorktreeDir verifies cwd is the worktree directory.
func TestRunHooks_WritesFileInWorktreeDir(t *testing.T) {
	dir := t.TempDir()
	captureStdout(t, func() {
		runHooks("myref", "myworktree", dir, "/fake/root", []string{"touch hook-sentinel.txt"})
	})
	if _, err := os.Stat(filepath.Join(dir, "hook-sentinel.txt")); os.IsNotExist(err) {
		t.Errorf("expected hook-sentinel.txt in worktree dir %s, but not found", dir)
	}
}
