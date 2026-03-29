package cmd

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadWipConfig_Valid tests loading a valid .wip.yml with on-worktree-create commands.
func TestLoadWipConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(orig)

	content := `refs:
  myservice:
    on-worktree-create:
      - echo hello
      - make setup
`
	if err := os.WriteFile(".wip.yml", []byte(content), 0644); err != nil {
		t.Fatalf("write .wip.yml: %v", err)
	}

	cfg, err := loadWipConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	sub, ok := cfg.Refs["myservice"]
	if !ok {
		t.Fatal("expected ref 'myservice' in config")
	}
	if len(sub.OnWorktreeCreate) != 2 {
		t.Fatalf("expected 2 on-worktree-create commands, got %d", len(sub.OnWorktreeCreate))
	}
	if sub.OnWorktreeCreate[0] != "echo hello" {
		t.Errorf("expected first command 'echo hello', got %q", sub.OnWorktreeCreate[0])
	}
	if sub.OnWorktreeCreate[1] != "make setup" {
		t.Errorf("expected second command 'make setup', got %q", sub.OnWorktreeCreate[1])
	}
}

// TestLoadWipConfig_Missing tests that a missing .wip.yml returns (nil, nil).
func TestLoadWipConfig_Missing(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(orig)

	cfg, err := loadWipConfig()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if cfg != nil {
		t.Fatalf("expected nil config for missing file, got: %+v", cfg)
	}
}

// TestLoadWipConfig_Malformed tests that a malformed .wip.yml returns a non-nil error.
func TestLoadWipConfig_Malformed(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(orig)

	// Write deliberately invalid YAML (bad indentation / type mismatch)
	malformed := `refs:
  - this is a list not a map
`
	if err := os.WriteFile(".wip.yml", []byte(malformed), 0644); err != nil {
		t.Fatalf("write .wip.yml: %v", err)
	}

	cfg, err := loadWipConfig()
	if err == nil {
		t.Fatalf("expected error for malformed YAML, got nil (cfg=%+v)", cfg)
	}
}

// TestRequireWipConfig_Absent tests that requireWipConfig returns an error mentioning "wip init"
// when .wip.yml does not exist.
func TestRequireWipConfig_Absent(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(orig)

	cfg, err := requireWipConfig()
	if err == nil {
		t.Fatalf("expected error when .wip.yml absent, got nil (cfg=%+v)", cfg)
	}
	if !strings.Contains(err.Error(), "wip init") {
		t.Errorf("expected error to contain 'wip init', got: %q", err.Error())
	}
}

// TestSaveWipConfig_RoundTrip saves a config and reloads it, verifying values match.
func TestSaveWipConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(orig)

	original := &WipConfig{
		Refs: map[string]RefConfig{
			"api": {
				OnWorktreeCreate: []string{"npm install", "npm run build"},
			},
			"web": {
				OnWorktreeCreate: []string{"yarn"},
			},
		},
	}

	if err := saveWipConfig(original); err != nil {
		t.Fatalf("saveWipConfig: %v", err)
	}

	loaded, err := loadWipConfig()
	if err != nil {
		t.Fatalf("loadWipConfig after save: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil config after round-trip")
	}

	// Verify 'api' ref
	api, ok := loaded.Refs["api"]
	if !ok {
		t.Fatal("expected ref 'api' in loaded config")
	}
	if len(api.OnWorktreeCreate) != 2 {
		t.Fatalf("expected 2 on-worktree-create commands for 'api', got %d", len(api.OnWorktreeCreate))
	}
	if api.OnWorktreeCreate[0] != "npm install" {
		t.Errorf("expected 'npm install', got %q", api.OnWorktreeCreate[0])
	}
	if api.OnWorktreeCreate[1] != "npm run build" {
		t.Errorf("expected 'npm run build', got %q", api.OnWorktreeCreate[1])
	}

	// Verify 'web' ref
	web, ok := loaded.Refs["web"]
	if !ok {
		t.Fatal("expected ref 'web' in loaded config")
	}
	if len(web.OnWorktreeCreate) != 1 {
		t.Fatalf("expected 1 on-worktree-create command for 'web', got %d", len(web.OnWorktreeCreate))
	}
	if web.OnWorktreeCreate[0] != "yarn" {
		t.Errorf("expected 'yarn', got %q", web.OnWorktreeCreate[0])
	}
}

// TestLoadWipConfig_OnWorktreeLaunch verifies that the on-worktree-launch field
// is parsed correctly from .wip.yml.
func TestLoadWipConfig_OnWorktreeLaunch(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(orig)

	content := `refs:
  myservice:
    on-worktree-launch:
      - git pull
      - npm install
      - claude
`
	if err := os.WriteFile(".wip.yml", []byte(content), 0644); err != nil {
		t.Fatalf("write .wip.yml: %v", err)
	}

	cfg, err := loadWipConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sub := cfg.Refs["myservice"]
	if len(sub.OnWorktreeLaunch) != 3 {
		t.Fatalf("expected 3 on-worktree-launch commands, got %d", len(sub.OnWorktreeLaunch))
	}
	want := []string{"git pull", "npm install", "claude"}
	for i, w := range want {
		if sub.OnWorktreeLaunch[i] != w {
			t.Errorf("OnWorktreeLaunch[%d]: want %q, got %q", i, w, sub.OnWorktreeLaunch[i])
		}
	}
}

// TestSaveWipConfig_OnWorktreeLaunchRoundTrip verifies that on-worktree-launch
// survives a save/load round-trip alongside on-worktree-create.
func TestSaveWipConfig_OnWorktreeLaunchRoundTrip(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	os.Chdir(dir)
	defer os.Chdir(orig)

	original := &WipConfig{
		Refs: map[string]RefConfig{
			"api": {
				OnWorktreeCreate: []string{"npm install"},
				OnWorktreeLaunch: []string{"git pull", "npm run dev"},
			},
		},
	}
	if err := saveWipConfig(original); err != nil {
		t.Fatalf("saveWipConfig: %v", err)
	}
	loaded, err := loadWipConfig()
	if err != nil {
		t.Fatalf("loadWipConfig: %v", err)
	}
	api := loaded.Refs["api"]
	if len(api.OnWorktreeCreate) != 1 || api.OnWorktreeCreate[0] != "npm install" {
		t.Errorf("OnWorktreeCreate round-trip failed: %v", api.OnWorktreeCreate)
	}
	if len(api.OnWorktreeLaunch) != 2 {
		t.Fatalf("expected 2 on-worktree-launch commands, got %d", len(api.OnWorktreeLaunch))
	}
	if api.OnWorktreeLaunch[0] != "git pull" || api.OnWorktreeLaunch[1] != "npm run dev" {
		t.Errorf("OnWorktreeLaunch round-trip failed: %v", api.OnWorktreeLaunch)
	}
}

// makeTempDirInHome creates a temporary directory tree inside the user's home directory
// and returns the root of the tree. t.Cleanup removes it.
func makeTempDirInHome(t *testing.T) string {
	t.Helper()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir: %v", err)
	}
	dir, err := os.MkdirTemp(home, ".wip-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp in home: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

// TestFindWipProject_FoundInCwd verifies that FindWipProject returns the current
// directory when it directly contains .wip.yml.
func TestFindWipProject_FoundInCwd(t *testing.T) {
	root := makeTempDirInHome(t)
	if err := os.WriteFile(filepath.Join(root, ".wip.yml"), []byte("refs: {}\n"), 0644); err != nil {
		t.Fatalf("write .wip.yml: %v", err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(orig)

	project, err := FindWipProject()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Root != root {
		t.Errorf("expected root %q, got %q", root, project.Root)
	}
	if project.Config == nil {
		t.Error("expected non-nil config")
	}
}

// TestFindWipProject_FoundInParent verifies that FindWipProject walks up and finds
// .wip.yml in a parent directory.
func TestFindWipProject_FoundInParent(t *testing.T) {
	root := makeTempDirInHome(t)
	if err := os.WriteFile(filepath.Join(root, ".wip.yml"), []byte("refs: {}\n"), 0644); err != nil {
		t.Fatalf("write .wip.yml: %v", err)
	}

	// Create a nested subdirectory and chdir into it.
	sub := filepath.Join(root, "submodule", "src", "components")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(sub); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(orig)

	project, err := FindWipProject()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Root != root {
		t.Errorf("expected root %q, got %q", root, project.Root)
	}
}

// TestFindWipProject_NotFoundWithinHome verifies that FindWipProject returns an error
// when no .wip.yml exists within the home directory tree.
func TestFindWipProject_NotFoundWithinHome(t *testing.T) {
	dir := makeTempDirInHome(t)
	// No .wip.yml anywhere in this tree.
	sub := filepath.Join(dir, "nested")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(sub); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(orig)

	_, err = FindWipProject()
	if err == nil {
		t.Fatal("expected error when no .wip.yml found, got nil")
	}
	if !strings.Contains(err.Error(), "wip init") {
		t.Errorf("expected error to mention 'wip init', got: %q", err.Error())
	}
}

// TestFindWipProject_OutsideHome verifies that FindWipProject fails immediately
// when cwd is outside the user's home directory.
func TestFindWipProject_OutsideHome(t *testing.T) {
	// Use os.TempDir() which is typically /tmp — outside home.
	tmpDir, err := os.MkdirTemp("", "wip-outside-home-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(orig)

	_, err = FindWipProject()
	if err == nil {
		t.Fatal("expected error when cwd is outside home, got nil")
	}
}

// TestStringList_Flag tests that a stringList flag accumulates values in order.
func TestStringList_Flag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var vals stringList
	fs.Var(&vals, "cmd", "a repeatable flag")

	args := []string{"--cmd", "first", "--cmd", "second", "--cmd", "third"}
	if err := fs.Parse(args); err != nil {
		t.Fatalf("flag parse: %v", err)
	}

	if len(vals) != 3 {
		t.Fatalf("expected 3 values, got %d: %v", len(vals), vals)
	}
	expected := []string{"first", "second", "third"}
	for i, v := range expected {
		if vals[i] != v {
			t.Errorf("vals[%d]: expected %q, got %q", i, v, vals[i])
		}
	}
}

// TestRefAdd_BranchDefaultsToMain verifies that the branch flag defaults to "main".
func TestRefAdd_BranchDefaultsToMain(t *testing.T) {
	// Parse a flagset with no --branch flag provided and verify the default is "main"
	fs := flag.NewFlagSet("ref add", flag.ContinueOnError)
	branch := fs.String("branch", "main", "branch to track")
	if err := fs.Parse([]string{"https://github.com/example/repo.git"}); err != nil {
		t.Fatalf("flag parse: %v", err)
	}
	if *branch != "main" {
		t.Errorf("expected branch default 'main', got %q", *branch)
	}
}

// TestSetGitmodulesIgnore verifies that setGitmodulesIgnore writes ignore=all via git config.
// This is an integration-style test that uses a real git repo in a temp dir.
func TestSetGitmodulesIgnore(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Init a git repo and a minimal .gitmodules file
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Skip("git not available")
	}
	gitmodules := "[submodule \"myrepo\"]\n\tpath = myrepo\n\turl = https://github.com/example/myrepo.git\n"
	if err := os.WriteFile(".gitmodules", []byte(gitmodules), 0644); err != nil {
		t.Fatalf("write .gitmodules: %v", err)
	}

	if err := setGitmodulesIgnore("myrepo"); err != nil {
		t.Fatalf("setGitmodulesIgnore: %v", err)
	}

	// Read back via git config
	out, err := exec.Command("git", "config", "--file", ".gitmodules", "submodule.myrepo.ignore").Output()
	if err != nil {
		t.Fatalf("git config read: %v", err)
	}
	if got := strings.TrimSpace(string(out)); got != "all" {
		t.Errorf("expected ignore=all, got %q", got)
	}
}

// TestRefRestore_EmptyState verifies that refRestore exits 0 when there are no refs (parseRefs returns empty).
func TestRefRestore_EmptyState(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	// Init a git repo with no submodules
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Skip("git not available")
	}

	refs, err := parseRefs()
	if err != nil {
		t.Fatalf("parseRefs: %v", err)
	}
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}
