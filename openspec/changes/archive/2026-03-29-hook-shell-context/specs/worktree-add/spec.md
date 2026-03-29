## MODIFIED Requirements

### Requirement: on-worktree-create hooks run after successful worktree creation
After `git worktree add` completes successfully, `wip worktree add` SHALL check `.wip.yml` for an `on-worktree-create` list under the target ref. If present, each command SHALL be executed via `sh -c` in list order with the working directory set to the newly created worktree path. Before each command runs, the standard `WIP_*` environment variables SHALL be injected into the subprocess environment.

#### Scenario: on-worktree-create hook executes successfully
- **WHEN** the ref has an `on-worktree-create` list in `.wip.yml` and all commands exit with code 0
- **THEN** all commands run in order inside the worktree directory and the CLI exits with code 0

#### Scenario: No on-worktree-create hook configured
- **WHEN** the ref has no `on-worktree-create` entry in `.wip.yml`
- **THEN** `wip worktree add` completes after git without running any hooks

#### Scenario: Ref has no entry in .wip.yml
- **WHEN** the ref name does not appear in the `.wip.yml` refs map
- **THEN** `wip worktree add` completes after git without running any hooks

#### Scenario: Hook uses WIP_REF_NAME env var
- **WHEN** a hook is configured as `"echo $WIP_REF_NAME"` and the ref is `backend`
- **THEN** the hook prints `backend`

#### Scenario: Hook uses WIP_WORKTREE_PATH env var
- **WHEN** a hook is configured as `"ls $WIP_WORKTREE_PATH"`
- **THEN** the hook lists the contents of the newly created worktree directory

#### Scenario: Hook uses shell compound command
- **WHEN** a hook is configured as `"npm install && npm run build"`
- **THEN** both commands execute in sequence inside the worktree directory

### Requirement: on-worktree-create hook failure produces a warning but leaves the worktree intact
If any command in the `on-worktree-create` list exits with a non-zero code, the CLI SHALL print a warning to stderr identifying which command failed. The CLI SHALL continue running any remaining commands in the list. The worktree directory SHALL NOT be removed. The CLI SHALL exit with code 0.

#### Scenario: One hook command fails
- **WHEN** a command in the `on-worktree-create` list exits non-zero
- **THEN** a warning is printed to stderr, subsequent commands still run, the worktree remains on disk, and the CLI exits with code 0

#### Scenario: All hook commands fail
- **WHEN** all commands in the `on-worktree-create` list exit non-zero
- **THEN** a warning is printed for each failure, the worktree remains on disk, and the CLI exits with code 0

### Requirement: on-worktree-create commands run with cwd set to the worktree directory
Each command in the `on-worktree-create` list SHALL be executed with its working directory set to the absolute path of the newly created worktree (`<repo root>/worktrees/<ref>/<worktree>/`).

#### Scenario: Working directory is the worktree
- **WHEN** an `on-worktree-create` command runs
- **THEN** its working directory is the newly created worktree directory, not the repo root
