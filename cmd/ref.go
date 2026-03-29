package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ref struct {
	name   string
	url    string
	branch string
}

// checkGitRepo reports an error if the current working directory is not a git repository.
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

// parseRefs reads all registered refs from .gitmodules (name, url, branch).
// Returns an empty slice (not an error) when no refs are registered.
func parseRefs() ([]ref, error) {
	cmd := exec.Command("git", "config", "--file", ".gitmodules", "--get-regexp", `submodule\..*\.url`)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read .gitmodules: %w", err)
	}

	var refs []ref
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
		url := parts[1]

		// Read branch for this ref; default to "main" if absent or on error.
		branch := "main"
		branchCmd := exec.Command("git", "config", "--file", ".gitmodules", fmt.Sprintf("submodule.%s.branch", name))
		if branchOut, err := branchCmd.Output(); err == nil {
			if b := strings.TrimSpace(string(branchOut)); b != "" {
				branch = b
			}
		}

		refs = append(refs, ref{name: name, url: url, branch: branch})
	}
	return refs, nil
}

// refExists reports whether a ref with the given name is registered in .gitmodules.
func refExists(name string) (bool, error) {
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

func Ref(args []string) {
	if len(args) < 1 {
		refUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		refAdd(args[1:])
	case "list":
		refList(args[1:])
	case "remove":
		refRemove(args[1:])
	case "sync":
		refSync(args[1:])
	case "restore":
		refRestore(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown ref command: %s\n", args[0])
		refUsage()
		os.Exit(1)
	}
}

func refUsage() {
	fmt.Fprintln(os.Stderr, "usage: wip ref <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  add <url>      add a git ref")
	fmt.Fprintln(os.Stderr, "  list           list registered refs")
	fmt.Fprintln(os.Stderr, "  remove <name>  fully remove a ref by name")
	fmt.Fprintln(os.Stderr, "  sync           update all refs concurrently")
	fmt.Fprintln(os.Stderr, "  restore        initialize all refs concurrently")
}
