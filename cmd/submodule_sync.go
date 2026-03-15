package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func submoduleSync(args []string) {
	fs := flag.NewFlagSet("submodule sync", flag.ExitOnError)
	name := fs.String("name", "", "sync only the named submodule (optional)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip submodule sync [--name <name>]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "flags:")
		fmt.Fprintln(os.Stderr, "  --name       sync only the named submodule (optional)")
	}
	fs.Parse(args)

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *name != "" {
		exists, err := submoduleExists(*name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !exists {
			fmt.Fprintf(os.Stderr, "submodule %q not found\n", *name)
			os.Exit(1)
		}

		out, err := exec.Command("git", "submodule", "update", "--init", "--remote", *name).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stdout, "✗ %s: %s\n", *name, string(out))
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "✓ %s\n", *name)
		return
	}

	subs, err := parseSubmodules()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(subs) == 0 {
		fmt.Println("no submodules to sync")
		os.Exit(0)
	}

	type syncResult struct {
		name string
		err  string
		ok   bool
	}

	results := make([]syncResult, len(subs))
	var wg sync.WaitGroup

	for i, sub := range subs {
		wg.Add(1)
		go func(idx int, subName string) {
			defer wg.Done()
			out, err := exec.Command("git", "submodule", "update", "--init", "--remote", subName).CombinedOutput()
			if err != nil {
				results[idx] = syncResult{name: subName, err: string(out), ok: false}
			} else {
				results[idx] = syncResult{name: subName, ok: true}
			}
		}(i, sub.name)
	}

	wg.Wait()

	anyFailed := false
	for _, r := range results {
		if r.ok {
			fmt.Printf("✓ %s\n", r.name)
		} else {
			fmt.Printf("✗ %s: %s\n", r.name, r.err)
			anyFailed = true
		}
	}

	if anyFailed {
		os.Exit(1)
	}
}
