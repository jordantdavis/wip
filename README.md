# wip

`wip` gives you fast, isolated workspaces across your repos — one command to spin up a branch, one to tear it down. Pull multiple independent repositories together under a single working directory so you can develop and navigate the full system in context.

---

## Quick Start

```bash
mkdir workspace && cd workspace
wip init
wip ref add https://github.com/org/backend.git
wip ref add https://github.com/org/frontend.git
wip worktree add backend my-feature
wip worktree add frontend my-feature
```

This gives you:

```
workspace/
├── backend/              ← ref checkout
├── frontend/             ← ref checkout
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

All `wip ref` and `wip worktree` commands work from any subdirectory of a wip project. `wip` walks up from the current directory to find the nearest `.wip.yml` and runs relative to that root automatically.

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

### `wip root`

Print the absolute path of the wip project root. Useful for scripting or navigating to the project root from a subdirectory.

```bash
wip root
# /Users/jordan/workspace
```

Works from any subdirectory:

```bash
cd workspace/backend/src/handlers
wip root
# /Users/jordan/workspace
```

---

### `wip ref add`

Add a git ref (external repo) to the workspace. Refs are git submodules configured to always track a branch HEAD and never dirty the parent repo's git status.

```
wip ref add [--name <name>] [--branch <branch>] [--on-worktree-create <cmd>] [--on-worktree-launch <cmd>] <url>
```

| Flag | Description |
|---|---|
| `--name` | Ref name and checkout directory. Defaults to the repo name from the URL. |
| `--branch` | Branch to track. Defaults to `main`. |
| `--on-worktree-create` | Shell command to run in a new worktree after creation. Repeatable. |
| `--on-worktree-launch` | Shell command to run when `wip worktree launch` is called for this ref. Repeatable. |

The URL must be one of: `https://`, `http://`, `git://`, or `git@<host>:<path>`.

```bash
# Basic add
wip ref add https://github.com/org/backend.git

# Override name and branch
wip ref add --name api --branch develop https://github.com/org/backend.git

# Register setup hooks (run in each new worktree)
wip ref add \
  --on-worktree-create "npm install" \
  --on-worktree-create "npm run build" \
  https://github.com/org/frontend.git

# Register launch hooks (run on wip worktree launch)
wip ref add \
  --on-worktree-launch "npm run dev" \
  https://github.com/org/frontend.git
```

---

### `wip ref list`

List all registered refs.

```bash
wip ref list
# backend   main  https://github.com/org/backend.git
# frontend  main  https://github.com/org/frontend.git
```

---

### `wip ref remove`

Fully remove a ref by name. Deinits it, removes the tracked path, and cleans up `.git/modules/<name>`.

```
wip ref remove <name>
```

```bash
wip ref remove backend
```

---

### `wip ref sync`

Update refs to the latest remote state. Syncs all refs in parallel by default.

```
wip ref sync [--name <name>]
```

| Flag | Description |
|---|---|
| `--name` | Sync only the named ref. |

```bash
# Sync all (parallel)
wip ref sync

# Sync one
wip ref sync --name backend
```

Output:
```
✓ backend
✓ frontend
```

---

### `wip ref restore`

Initialize and sync all registered refs. Use this after cloning the workspace repo to a new machine.

```bash
wip ref restore
```

Output:
```
✓ backend
✓ frontend
```

Always fetches the current branch HEAD (`--remote`) so teammates get the same view of the codebase rather than a stale committed SHA.

---

### `wip worktree add`

Create a new worktree in a ref. Creates a new branch by default.

```
wip worktree add [--existing-branch] <ref> <worktree>
```

| Flag | Description |
|---|---|
| `--existing-branch` | Check out an existing branch instead of creating a new one. |

Worktree names must match `[a-zA-Z0-9_-]+`. The name is used as both the directory name and the branch name.

The worktree is created at `worktrees/<ref>/<worktree>/`. Any `on-worktree-create` hooks registered for the ref run automatically in the new worktree directory.

```bash
# New branch
wip worktree add backend my-feature

# Existing branch
wip worktree add --existing-branch backend main
```

---

### `wip worktree list`

List all worktrees across all refs.

```bash
wip worktree list
# REF       WORKTREE    BRANCH
# backend   my-feature  my-feature
# frontend  my-feature  my-feature
```

---

### `wip worktree launch`

Run `on-worktree-launch` hooks for an existing worktree. Use this to start services or open editors associated with a worktree.

```
wip worktree launch <ref> <worktree>
```

```bash
wip worktree launch backend my-feature
# ✓ npm run dev
```

If no `on-worktree-launch` hooks are configured for the ref, the command exits with a message.

---

### `wip worktree remove`

Remove a worktree from a ref.

```
wip worktree remove [--delete-branch] <ref> <worktree>
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
refs:
  backend:
    on-worktree-create:
      - go mod download
    on-worktree-launch:
      - go run .
  frontend:
    on-worktree-create:
      - npm install
      - npm run build
    on-worktree-launch:
      - npm run dev
```

`on-worktree-create` commands are executed in order inside the new worktree directory whenever `wip worktree add` runs for that ref.

`on-worktree-launch` commands are executed in order inside the worktree directory whenever `wip worktree launch` runs for that ref.

Each command is reported as a success or failure:

```
✓ npm install
✓ npm run build
```
