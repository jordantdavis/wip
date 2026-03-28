package main

import (
	"fmt"
	"os"

	"github.com/jordantdavis/wip/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cmd.Init(os.Args[2:])
	case "version":
		cmd.Version()
	case "root":
		cmd.Root(os.Args[2:])
	case "submodule", "worktree":
		project, err := cmd.FindWipProject()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := os.Chdir(project.Root); err != nil {
			fmt.Fprintf(os.Stderr, "failed to change to project root: %v\n", err)
			os.Exit(1)
		}
		switch os.Args[1] {
		case "submodule":
			cmd.Submodule(os.Args[2:])
		case "worktree":
			cmd.Worktree(os.Args[2:])
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: wip <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  init        initialize a wip project in the current directory")
	fmt.Fprintln(os.Stderr, "  root        print the project root directory")
	fmt.Fprintln(os.Stderr, "  submodule   manage git submodules")
	fmt.Fprintln(os.Stderr, "  worktree    manage git worktrees")
	fmt.Fprintln(os.Stderr, "  version     print version information")
}
