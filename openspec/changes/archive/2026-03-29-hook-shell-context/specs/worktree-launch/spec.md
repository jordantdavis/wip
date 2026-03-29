## MODIFIED Requirements

### Requirement: wip worktree launch runs on-worktree-launch hooks sequentially in the worktree directory
`wip worktree launch` SHALL execute each command in the `on-worktree-launch` list via `sh -c` in order, with the working directory set to the worktree path. Before each command runs, the standard `WIP_*` environment variables SHALL be injected into the subprocess environment. Commands run blocking and sequentially. The final command MAY be interactive and take over the terminal.

#### Scenario: Single hook runs successfully
- **WHEN** the ref has one `on-worktree-launch` command and it exits with code 0
- **THEN** the command runs in the worktree directory, ✓ is printed, and the CLI exits with code 0

#### Scenario: Multiple hooks run in order
- **WHEN** the ref has multiple `on-worktree-launch` commands
- **THEN** each command runs sequentially in the worktree directory in list order

#### Scenario: Hook failure prints warning and continues
- **WHEN** a command in the `on-worktree-launch` list exits non-zero
- **THEN** ✗ is printed for that command, subsequent commands still run, and the CLI exits with code 0

#### Scenario: Last command is interactive
- **WHEN** the last `on-worktree-launch` command is an interactive process (e.g. `claude`)
- **THEN** it takes over the terminal and the CLI exits when that process exits

#### Scenario: Hook uses WIP_WORKTREE_NAME env var
- **WHEN** a hook is configured as `"echo $WIP_WORKTREE_NAME"` and the worktree is `my-feature`
- **THEN** the hook prints `my-feature`

#### Scenario: Hook uses WIP_WORKTREE_PATH to open editor
- **WHEN** a hook is configured as `"open -n /Applications/Cursor.app --args \"$WIP_WORKTREE_PATH\""`
- **THEN** the editor opens at the worktree path
