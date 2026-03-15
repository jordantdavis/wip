## Why

The project currently documents build commands only in CLAUDE.md, with no standard entry point for running them. Adding a Makefile provides a consistent, discoverable interface for local development and will serve as the foundation for GitHub Actions CI workflows.

## What Changes

- Add a `Makefile` at the repo root with targets for building, testing, formatting, format-checking, vetting, installing, cleaning, and local verification

## Capabilities

### New Capabilities

- `makefile`: A Makefile providing standard targets for the full local dev and CI workflow

### Modified Capabilities

## Impact

- Adds `Makefile` to the repo root
- No changes to existing Go source code
- CLAUDE.md commands section may be updated to reference `make` targets
