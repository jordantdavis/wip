package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestValidateBranchName_ValidSimple verifies a simple branch name is accepted.
// This is an integration test that invokes git.
func TestValidateBranchName_ValidSimple(t *testing.T) {
	if err := validateBranchName("my-feature"); err != nil {
		t.Errorf("expected nil error for valid name 'my-feature', got: %v", err)
	}
}

// TestValidateBranchName_ValidSlash verifies a slash-delimited branch name is accepted.
func TestValidateBranchName_ValidSlash(t *testing.T) {
	if err := validateBranchName("feature/my-thing"); err != nil {
		t.Errorf("expected nil error for valid name 'feature/my-thing', got: %v", err)
	}
}

// TestValidateBranchName_Invalid verifies an invalid branch name returns a non-nil error.
func TestValidateBranchName_Invalid(t *testing.T) {
	if err := validateBranchName("..bad"); err == nil {
		t.Error("expected non-nil error for invalid name '..bad', got nil")
	}
}

// TestWorktreePathSegment_NoSlash verifies a name with no slashes is returned unchanged.
func TestWorktreePathSegment_NoSlash(t *testing.T) {
	got := worktreePathSegment("my-feature")
	want := "my-feature"
	if got != want {
		t.Errorf("worktreePathSegment(%q) = %q, want %q", "my-feature", got, want)
	}
}

// TestWorktreePathSegment_SingleSlash verifies a single slash is replaced with a hyphen.
func TestWorktreePathSegment_SingleSlash(t *testing.T) {
	got := worktreePathSegment("feature/my-thing")
	want := "feature-my-thing"
	if got != want {
		t.Errorf("worktreePathSegment(%q) = %q, want %q", "feature/my-thing", got, want)
	}
}

// TestWorktreePathSegment_MultipleSlashes verifies multiple slashes are all replaced with hyphens.
func TestWorktreePathSegment_MultipleSlashes(t *testing.T) {
	got := worktreePathSegment("team/user/ticket-123")
	want := "team-user-ticket-123"
	if got != want {
		t.Errorf("worktreePathSegment(%q) = %q, want %q", "team/user/ticket-123", got, want)
	}
}

// TestWorktreeList_ThreeColumnFormat verifies the three-column output format produced
// by worktreeList. It sets up a fake worktrees directory with a real git repo in the
// leaf directory so that `git branch --show-current` returns a known value, then
// captures stdout while running worktreeList.
func TestWorktreeList_ThreeColumnFormat(t *testing.T) {
	root := t.TempDir()

	// Create worktrees/myservice/feature-my-thing/
	submodule := "myservice"
	pathSeg := "feature-my-thing"
	worktreeDir := filepath.Join(root, "worktrees", submodule, pathSeg)
	if err := os.MkdirAll(worktreeDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Init a git repo in the worktree directory on branch "feature/my-thing".
	gitRun := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	gitRun(worktreeDir, "init", "-b", "feature/my-thing")
	gitRun(worktreeDir, "config", "user.email", "test@example.com")
	gitRun(worktreeDir, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(worktreeDir, ".keep"), []byte(""), 0644); err != nil {
		t.Fatalf("write .keep: %v", err)
	}
	gitRun(worktreeDir, "add", ".")
	gitRun(worktreeDir, "commit", "-m", "init")

	// Create a .git directory in root so checkGitRepo() passes.
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	// Change working directory to root so repoRoot() returns root.
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(orig)

	// Capture stdout.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	worktreeList([]string{})

	w.Close()
	os.Stdout = oldStdout

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

	output := strings.TrimSpace(buf.String())
	t.Logf("worktreeList output: %q", output)

	// Expect a header line followed by one data line.
	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines (header + data), got %d: %v", len(lines), lines)
	}

	// The format is "%s  %s  %s\n" — split on runs of whitespace.
	parts := strings.Fields(lines[1])
	if len(parts) != 3 {
		t.Fatalf("expected 3 fields in output line, got %d: %q", len(parts), lines[0])
	}

	wantSubmodule := submodule
	wantPathSeg := pathSeg
	wantBranch := "feature/my-thing"

	if parts[0] != wantSubmodule {
		t.Errorf("column 1 (submodule): want %q, got %q", wantSubmodule, parts[0])
	}
	if parts[1] != wantPathSeg {
		t.Errorf("column 2 (path-segment): want %q, got %q", wantPathSeg, parts[1])
	}
	if parts[2] != wantBranch {
		t.Errorf("column 3 (branch): want %q, got %q", wantBranch, parts[2])
	}

	// Verify the raw data line uses double-space separators as specified.
	wantLine := fmt.Sprintf("%s  %s  %s", wantSubmodule, wantPathSeg, wantBranch)
	if lines[1] != wantLine {
		t.Errorf("output line format: want %q, got %q", wantLine, lines[1])
	}
}
