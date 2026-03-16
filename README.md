# wip

`wip` gives you fast, isolated workspaces across your repos — one command to spin up a branch, one to tear it down. Pull multiple independent repositories together under a single working directory so you can develop and navigate the full system in context.

---

## Quick Start

```bash
mkdir workspace && cd workspace
wip init
wip submodule add https://github.com/org/backend.git
wip submodule add https://github.com/org/frontend.git
wip worktree add backend my-feature
wip worktree add frontend my-feature
```

This gives you:

```
workspace/
├── backend/              ← submodule checkout
├── frontend/             ← submodule checkout
└── worktrees/
    ├── backend/
    │   └── my-feature/   ← isolated branch checkout
    └── frontend/
        └── my-feature/   ← isolated branch checkout
```

---

## Installation

### GitHub Releases

Download a pre-built binary from [Releases](https://github.com/jordantdavis/wip/releases) and place it on your `PATH`.

### Go Install

```bash
go install github.com/jordantdavis/wip@<version>
```

### Build from Source

```bash
git clone https://github.com/jordantdavis/wip.git
cd wip
go build -o wip .
```

---

## Commands

### `wip init`

Initialize a `wip` workspace in the current directory.

Creates a `.wip.yml` config file. If the current directory is not already a git repository, runs `git init` first.

```bash
wip init
```

---

### `wip version`

Print version and platform information.

```bash
wip version
# wip v0.0.1 (darwin/arm64)
```

---

### `wip submodule add`

Add a git submodule to the workspace.

```
wip submodule add [--name <name>] [--on-worktree-create <cmd>] <url>
```

| Flag | Description |
|---|---|
| `--name` | Submodule name and checkout directory. Defaults to the repo name from the URL. |
| `--on-worktree-create` | Shell command to run in a new worktree after creation. Repeatable. |

The URL must be one of: `https://`, `http://`, `git://`, or `git@<host>:<path>`.

```bash
# Basic add
wip submodule add https://github.com/org/backend.git

# Override name
wip submodule add --name api https://github.com/org/backend.git

# Register setup hooks (run in each new worktree)
wip submodule add \
  --on-worktree-create "npm install" \
  --on-worktree-create "npm run build" \
  https://github.com/org/frontend.git
```

---

### `wip submodule list`

List all registered submodules.

```bash
wip submodule list
# backend   https://github.com/org/backend.git
# frontend  https://github.com/org/frontend.git
```

---

### `wip submodule remove`

Fully remove a submodule by name. Deinits it, removes the tracked path, and cleans up `.git/modules/<name>`.

```
wip submodule remove <name>
```

```bash
wip submodule remove backend
```

---

### `wip submodule sync`

Update submodules to the latest remote state. Syncs all submodules in parallel by default.

```
wip submodule sync [--name <name>]
```

| Flag | Description |
|---|---|
| `--name` | Sync only the named submodule. |

```bash
# Sync all (parallel)
wip submodule sync

# Sync one
wip submodule sync --name backend
```

Output:
```
✓ backend
✓ frontend
```

---

### `wip worktree add`

Create a new worktree in a submodule. Creates a new branch by default.

```
wip worktree add [--existing-branch] <submodule> <worktree>
```

| Flag | Description |
|---|---|
| `--existing-branch` | Check out an existing branch instead of creating a new one. |

Worktree names must match `[a-zA-Z0-9_-]+`. The name is used as both the directory name and the branch name.

The worktree is created at `worktrees/<submodule>/<worktree>/`. Any `on-worktree-create` hooks registered for the submodule run automatically in the new worktree directory.

```bash
# New branch
wip worktree add backend my-feature

# Existing branch
wip worktree add --existing-branch backend main
```

---

### `wip worktree list`

List all worktrees across all submodules.

```bash
wip worktree list
# backend   my-feature
# frontend  my-feature
```

---

### `wip worktree remove`

Remove a worktree from a submodule.

```
wip worktree remove [--delete-branch] <submodule> <worktree>
```

| Flag | Description |
|---|---|
| `--delete-branch` | Also delete the associated branch. |

```bash
# Remove worktree only
wip worktree remove backend my-feature

# Remove worktree and branch
wip worktree remove --delete-branch backend my-feature
```

---

## Configuration

`wip init` creates a `.wip.yml` at the workspace root. You can also edit it directly.

```yaml
submodules:
  backend:
    on-worktree-create:
      - go mod download
  frontend:
    on-worktree-create:
      - npm install
      - npm run build
```

`on-worktree-create` commands are executed in order inside the new worktree directory whenever `wip worktree add` runs for that submodule. Each command is reported as a success or failure:

```
✓ npm install
✓ npm run build
```
