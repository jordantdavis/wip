package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func submoduleRemove(args []string) {
	fs := flag.NewFlagSet("submodule remove", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip submodule remove <name>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "args:")
		fmt.Fprintln(os.Stderr, "  name         name of the submodule to remove")
	}
	fs.Parse(args)

	positional := fs.Args()
	if len(positional) < 1 {
		fs.Usage()
		os.Exit(1)
	}

	name := positional[0]

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	exists, err := submoduleExists(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !exists {
		fmt.Fprintf(os.Stderr, "submodule %q not found\n", name)
		os.Exit(1)
	}

	// Step a: git submodule deinit -f <name>
	deinit := exec.Command("git", "submodule", "deinit", "-f", name)
	deinit.Stdout = os.Stdout
	deinit.Stderr = os.Stderr
	if err := deinit.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Fprintln(os.Stderr, "deinit failed")
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "deinit failed")
		os.Exit(1)
	}

	// Step b: git rm -f <name>
	gitRm := exec.Command("git", "rm", "-f", name)
	gitRm.Stdout = os.Stdout
	gitRm.Stderr = os.Stderr
	if err := gitRm.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Fprintln(os.Stderr, "git rm failed")
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "git rm failed")
		os.Exit(1)
	}

	// Step c: remove .git/modules/<name>
	if err := os.RemoveAll(filepath.Join(".git", "modules", name)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
