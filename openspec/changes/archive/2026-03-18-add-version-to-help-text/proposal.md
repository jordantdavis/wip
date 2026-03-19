## Why

The `version` command is already implemented and routed in `main.go`, but it doesn't appear in the CLI's help output. Users have no way to discover it without reading the source.

## What Changes

- Add `version` to the `printUsage()` output in `main.go`

## Capabilities

### New Capabilities
- `version-discoverability`: The `version` command is listed in CLI help output alongside other commands

### Modified Capabilities

## Impact

- `main.go`: `printUsage()` function gets one additional line
