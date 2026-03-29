package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type syncRefResult struct {
	name string
	err  string
	ok   bool
}

func refSync(args []string) {
	fs := flag.NewFlagSet("ref sync", flag.ExitOnError)
	name := fs.String("name", "", "sync only the named ref")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip ref sync [--name <name>]")
	}
	fs.Parse(args)

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *name != "" {
		// Single ref sync
		exists, err := refExists(*name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !exists {
			fmt.Fprintf(os.Stderr, "ref %q not found\n", *name)
			os.Exit(1)
		}

		out, err := exec.Command("git", "submodule", "update", "--remote", *name).CombinedOutput()
		if err != nil {
			fmt.Printf("\u2717 %s: %s\n", *name, string(out))
			os.Exit(1)
		}
		fmt.Printf("\u2713 %s\n", *name)
		return
	}

	// Sync all refs concurrently
	refs, err := parseRefs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(refs) == 0 {
		fmt.Println("no refs to sync")
		os.Exit(0)
	}

	results := make([]syncRefResult, len(refs))
	var wg sync.WaitGroup
	for i, r := range refs {
		wg.Add(1)
		go func(idx int, refName string) {
			defer wg.Done()
			out, err := exec.Command("git", "submodule", "update", "--remote", refName).CombinedOutput()
			if err != nil {
				results[idx] = syncRefResult{name: refName, err: string(out), ok: false}
			} else {
				results[idx] = syncRefResult{name: refName, ok: true}
			}
		}(i, r.name)
	}
	wg.Wait()

	anyFailed := false
	for _, res := range results {
		if res.ok {
			fmt.Printf("\u2713 %s\n", res.name)
		} else {
			fmt.Printf("\u2717 %s: %s\n", res.name, res.err)
			anyFailed = true
		}
	}

	if anyFailed {
		os.Exit(1)
	}
}
