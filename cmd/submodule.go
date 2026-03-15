package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type submodule struct {
	name string
	url  string
}

// parseSubmodules reads all registered submodules from .gitmodules.
// Returns an empty slice (not an error) when no submodules are registered.
func parseSubmodules() ([]submodule, error) {
	cmd := exec.Command("git", "config", "--file", ".gitmodules", "--get-regexp", `submodule\..*\.url`)
	out, err := cmd.Output()
	if err != nil {
		// exit code 1 means the key pattern matched nothing (no submodules)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read .gitmodules: %w", err)
	}

	var subs []submodule
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		// key is "submodule.<name>.url"; name may itself contain dots
		key := parts[0]
		keyParts := strings.Split(key, ".")
		if len(keyParts) < 3 {
			continue
		}
		name := strings.Join(keyParts[1:len(keyParts)-1], ".")
		subs = append(subs, submodule{name: name, url: parts[1]})
	}
	return subs, nil
}

// submoduleExists reports whether a submodule with the given name is registered in .gitmodules.
func submoduleExists(name string) (bool, error) {
	cmd := exec.Command("git", "config", "--file", ".gitmodules", "--get", fmt.Sprintf("submodule.%s.url", name))
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("failed to check .gitmodules: %w", err)
	}
	return true, nil
}

func Submodule(args []string) {
	if len(args) < 1 {
		submoduleUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		submoduleAdd(args[1:])
	case "list":
		submoduleList(args[1:])
	case "remove":
		submoduleRemove(args[1:])
	case "sync":
		submoduleSync(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown submodule command: %s\n", args[0])
		submoduleUsage()
		os.Exit(1)
	}
}

func submoduleUsage() {
	fmt.Fprintln(os.Stderr, "usage: wip submodule <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  add <url>      add a git submodule")
	fmt.Fprintln(os.Stderr, "  list           list registered submodules")
	fmt.Fprintln(os.Stderr, "  remove <name>  fully remove a submodule by name")
	fmt.Fprintln(os.Stderr, "  sync           update all submodules concurrently")
}

func checkGitRepo() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	info, err := os.Stat(filepath.Join(cwd, ".git"))
	if err != nil || !info.IsDir() {
		return errors.New("not a git repository (no .git directory found)")
	}
	return nil
}
