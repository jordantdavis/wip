package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var version = "unknown"

func getVersion() string {
	if version != "unknown" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "unknown"
}

func Version() {
	fmt.Printf("wip %s (%s/%s)\n", getVersion(), runtime.GOOS, runtime.GOARCH)
}
