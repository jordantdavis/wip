## Why

`wip` subcommands currently only work when invoked from the project root (the directory containing `.wip.yml`). This forces users to `cd` to the project root before running any `wip` command, which breaks the natural workflow of operating from within a submodule or nested directory.

## What Changes

- `wip` will walk upward from the current directory (bounded by the user's home directory) to locate the nearest `.wip.yml` and treat that directory as the project root
- All subcommands (`submodule`, `worktree`) will automatically run relative to the discovered project root
- New `wip root` command prints the absolute path of the discovered project root
- If invoked from outside the home directory, subcommands fail immediately with a clear error
- If no `.wip.yml` is found within the home directory tree, subcommands fail with the existing "run wip init" message
- `wip init` and `wip version` are unaffected — they do not require an existing project
- README updated to document subdirectory-aware behavior and `wip root`

## Capabilities

### New Capabilities

- `wip-project-discovery`: Locating the wip project root by walking up the directory tree from the current working directory, bounded by the user's home directory
- `wip-root`: A command that prints the absolute path of the discovered project root

### Modified Capabilities

- `wip-config`: The config loading and saving functions must resolve paths relative to the discovered project root, not the current working directory

## Impact

- `cmd/config.go`: New `WipProject` type and `findWipProject()` function; `loadWipConfig` and `saveWipConfig` updated to take an absolute path
- `cmd/root.go`: New file implementing `wip root`
- `main.go`: Calls `findWipProject()` and `os.Chdir(project.Root)` before dispatching non-init/version subcommands; adds `root` case to dispatch
- All subcommands: Unchanged — they continue to rely on `os.Getwd()` which now points to the project root after the chdir
- `cmd/init.go`: Unchanged — init is intentionally root-only and does not participate in project discovery
- `README.md`: Documents subdirectory-aware behavior and `wip root` command
