package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type restoreRefResult struct {
	name string
	err  string
	ok   bool
}

func refRestore(args []string) {
	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	refs, err := parseRefs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(refs) == 0 {
		fmt.Println("no refs registered")
		os.Exit(0)
	}

	results := make([]restoreRefResult, len(refs))
	var wg sync.WaitGroup
	for i, r := range refs {
		wg.Add(1)
		go func(idx int, refName string) {
			defer wg.Done()
			out, err := exec.Command("git", "submodule", "update", "--init", "--remote", refName).CombinedOutput()
			if err != nil {
				results[idx] = restoreRefResult{name: refName, err: string(out), ok: false}
			} else {
				results[idx] = restoreRefResult{name: refName, ok: true}
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
