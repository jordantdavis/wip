package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func submoduleAdd(args []string) {
	fs := flag.NewFlagSet("submodule add", flag.ExitOnError)
	name := fs.String("name", "", "submodule name and checkout directory at repo root (optional)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip submodule add [--name <name>] <url>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "args:")
		fmt.Fprintln(os.Stderr, "  url          git remote URL (https://, http://, git://, or git@)")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "flags:")
		fmt.Fprintln(os.Stderr, "  --name       submodule name and checkout directory at repo root (optional)")
	}
	fs.Parse(args)

	positional := fs.Args()
	if len(positional) < 1 {
		fs.Usage()
		os.Exit(1)
	}

	url := positional[0]

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := validateURL(url); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *name != "" {
		if err := validateName(*name); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if err := runGitSubmoduleAdd(url, *name); err != nil {
		os.Exit(1)
	}
}

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
	// SSH form: git@<host>:<path>
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

func runGitSubmoduleAdd(url, name string) error {
	gitArgs := []string{"submodule", "add"}
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
