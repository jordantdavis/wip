package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func validateURL(url string) error {
	if url == "" {
		return errors.New("url must not be empty")
	}
	validPrefixes := []string{"https://", "http://", "git://"}
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(url, prefix) {
			return nil
		}
	}
	if strings.HasPrefix(url, "git@") && strings.Contains(url, ":") {
		return nil
	}
	return fmt.Errorf("invalid git remote URL: %q (expected https://, http://, git://, or git@<host>:<path>)", url)
}

func validateName(name string) error {
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("name %q must not contain path separators", name)
	}
	return nil
}

func refAdd(args []string) {
	fs := flag.NewFlagSet("ref add", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip ref add [--name <name>] [--branch <branch>] [--on-worktree-create <cmd>] [--on-worktree-launch <cmd>] <url>")
		fs.PrintDefaults()
	}

	name := fs.String("name", "", "ref name and checkout directory at repo root (optional)")
	branch := fs.String("branch", "main", "branch to track")
	var onWorktreeCreate stringList
	var onWorktreeLaunch stringList
	fs.Var(&onWorktreeCreate, "on-worktree-create", "command to run when worktree created (repeatable)")
	fs.Var(&onWorktreeLaunch, "on-worktree-launch", "command to run when worktree launched (repeatable)")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: url is required")
		fs.Usage()
		os.Exit(1)
	}
	url := fs.Arg(0)

	cfg, err := requireWipConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if err := validateURL(url); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *name != "" {
		if err := validateName(*name); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}

	if err := runRefAdd(url, *name, *branch); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	effectiveName := *name
	if effectiveName == "" {
		effectiveName = strings.TrimSuffix(path.Base(url), ".git")
	}

	if err := setGitmodulesIgnore(effectiveName); err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to set ignore=all in .gitmodules:", err)
	}

	if cfg.Refs == nil {
		cfg.Refs = make(map[string]RefConfig)
	}
	refCfg := cfg.Refs[effectiveName]
	refCfg.URL = url
	refCfg.Branch = *branch
	if len(onWorktreeCreate) > 0 {
		refCfg.OnWorktreeCreate = []string(onWorktreeCreate)
	}
	if len(onWorktreeLaunch) > 0 {
		refCfg.OnWorktreeLaunch = []string(onWorktreeLaunch)
	}
	cfg.Refs[effectiveName] = refCfg
	if err := saveWipConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to save .wip.yml:", err)
	}
}

func runRefAdd(url, name, branch string) error {
	gitArgs := []string{"submodule", "add", "-b", branch}
	if name != "" {
		gitArgs = append(gitArgs, "--name", name)
	}
	gitArgs = append(gitArgs, url)
	if name != "" {
		gitArgs = append(gitArgs, name)
	}
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	return nil
}

func setGitmodulesIgnore(name string) error {
	key := fmt.Sprintf("submodule.%s.ignore", name)
	cmd := exec.Command("git", "config", "--file", ".gitmodules", key, "all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
