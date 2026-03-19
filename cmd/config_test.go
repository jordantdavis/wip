package cmd

import (
	"flag"
	"os"
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

	content := `submodules:
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
	sub, ok := cfg.Submodules["myservice"]
	if !ok {
		t.Fatal("expected submodule 'myservice' in config")
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
	malformed := `submodules:
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
		Submodules: map[string]SubmoduleConfig{
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

	// Verify 'api' submodule
	api, ok := loaded.Submodules["api"]
	if !ok {
		t.Fatal("expected submodule 'api' in loaded config")
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

	// Verify 'web' submodule
	web, ok := loaded.Submodules["web"]
	if !ok {
		t.Fatal("expected submodule 'web' in loaded config")
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

	content := `submodules:
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
	sub := cfg.Submodules["myservice"]
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
		Submodules: map[string]SubmoduleConfig{
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
	api := loaded.Submodules["api"]
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
