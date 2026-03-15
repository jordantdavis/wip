## Why

When developing in a monorepo with multiple submodules, each new worktree requires manual environment setup (e.g., `npm install`, copying `.env` files) before it is usable. There is no mechanism to automate this per-submodule initialization today.

## What Changes

- `wip init` scaffolds a `.wip.yml` file at the repo root in addition to running `git init`
- `.wip.yml` is required for `wip submodule add` and `wip worktree add`; both commands fail with a nudge to `wip init` if it is absent
- `wip submodule add` gains a repeatable `--on-worktree-create` flag (placed before the URL, consistent with `--name`) that accepts ordered commands; these are written to `.wip.yml` under the submodule entry
- `wip worktree add` reads `.wip.yml` and, if an `on-worktree-create` list exists for the target submodule, runs each command in order inside the newly created worktree directory; hook failures print a warning but leave the worktree intact
- `.wip.yml` uses kebab-case keys and is processed with `go.yaml.in/yaml/v4`

## Capabilities

### New Capabilities

- `wip-config`: `.wip.yml` file format, location, lifecycle, and the shared validation helper that enforces its presence
- `worktree-lifecycle`: The `on-worktree-create` hook model — configuration shape, execution contract (cwd, ordering, failure behavior), and how hooks are stored and retrieved

### Modified Capabilities

- `wip-init`: Must now also scaffold an empty `.wip.yml` at the repo root
- `submodule-add`: Gains the repeatable `--on-worktree-create` flag and writes hooks to `.wip.yml`
- `worktree-add`: Must validate `.wip.yml` presence and execute `on-worktree-create` hooks after worktree creation

## Impact

- New dependency: `go.yaml.in/yaml/v4`
- New file: `.wip.yml` at repo root (created by `wip init`, required by add commands)
- New shared helper in `cmd/` for `.wip.yml` presence validation
- New config package or file for reading/writing `.wip.yml`
- `wip submodule add` and `wip worktree add` are the only commands affected by the `.wip.yml` requirement for now
