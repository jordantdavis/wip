package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func submoduleList(args []string) {
	fs := flag.NewFlagSet("submodule list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip submodule list")
	}
	fs.Parse(args)

	if err := checkGitRepo(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	subs, err := parseSubmodules()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(subs) == 0 {
		fmt.Println("no submodules found")
		os.Exit(0)
	}

	sort.Slice(subs, func(i, j int) bool {
		return subs[i].name < subs[j].name
	})

	for _, s := range subs {
		fmt.Printf("%s  %s\n", s.name, s.url)
	}

	os.Exit(0)
}
