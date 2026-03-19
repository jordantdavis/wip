## ADDED Requirements

### Requirement: --on-worktree-launch flag accepts ordered hook commands
`wip submodule add` SHALL accept a repeatable `--on-worktree-launch` flag placed before the URL. Each use appends a command string to the list in CLI argument order. The flag MAY be omitted entirely. When provided, the collected commands SHALL be written to `.wip.yml` under `submodules.<name>.on-worktree-launch` as an ordered list.

#### Scenario: Single --on-worktree-launch command
- **WHEN** the user runs `wip submodule add --on-worktree-launch "npm run dev" <url>`
- **THEN** `.wip.yml` is updated with `on-worktree-launch: ["npm run dev"]` under the submodule entry

#### Scenario: Multiple --on-worktree-launch commands preserve order
- **WHEN** the user runs `wip submodule add --on-worktree-launch "git pull" --on-worktree-launch "npm install" --on-worktree-launch "claude" <url>`
- **THEN** `.wip.yml` is updated with `on-worktree-launch: ["git pull", "npm install", "claude"]` in that order

#### Scenario: --on-worktree-launch omitted
- **WHEN** the user runs `wip submodule add <url>` without `--on-worktree-launch`
- **THEN** no `on-worktree-launch` entry is written to `.wip.yml` for this submodule

#### Scenario: Both --on-worktree-create and --on-worktree-launch provided
- **WHEN** the user provides both `--on-worktree-create` and `--on-worktree-launch` flags
- **THEN** both lists are written to `.wip.yml` under the same submodule entry independently
