package cmd

import (
	"fmt"
	"os"
)

// Root implements the `wip root` command. It prints the absolute path of the
// discovered wip project root to stdout and exits 0, or prints an error to
// stderr and exits non-zero if no project is found.
func Root(_ []string) {
	project, err := FindWipProject()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(project.Root)
}
