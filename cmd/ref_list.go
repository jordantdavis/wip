package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func refList(args []string) {
	fs := flag.NewFlagSet("ref list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: wip ref list")
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	refs, err := parseRefs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if len(refs) == 0 {
		fmt.Println("no refs found")
		return
	}

	sort.Slice(refs, func(i, j int) bool {
		return refs[i].name < refs[j].name
	})

	for _, r := range refs {
		fmt.Printf("%s  %s  %s\n", r.name, r.branch, r.url)
	}
}
