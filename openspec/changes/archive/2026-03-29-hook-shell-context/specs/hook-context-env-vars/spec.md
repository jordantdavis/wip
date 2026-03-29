## ADDED Requirements

### Requirement: Standard WIP_* env vars are injected for every hook execution
Before executing any hook command, the CLI SHALL inject the following environment variables into the hook subprocess environment. These variables SHALL be present for all current and future hook types.

| Variable | Description |
|---|---|
| `WIP_REF_NAME` | The ref name as defined in `.wip.yml` |
| `WIP_WORKTREE_NAME` | The worktree name as provided by the user |
| `WIP_WORKTREE_PATH` | Absolute path to the worktree directory |
| `WIP_ROOT` | Absolute path to the repository root |

#### Scenario: Env vars are available inside the hook
- **WHEN** a hook is configured as `"echo $WIP_REF_NAME $WIP_WORKTREE_NAME"`
- **THEN** the output contains the ref name and worktree name

#### Scenario: WIP_WORKTREE_PATH is the absolute path
- **WHEN** a hook reads `$WIP_WORKTREE_PATH`
- **THEN** the value is an absolute filesystem path to the worktree directory

#### Scenario: WIP_ROOT is the repo root
- **WHEN** a hook reads `$WIP_ROOT`
- **THEN** the value is the absolute path to the repository root containing `.git`

### Requirement: Hook env vars are inherited by subprocesses
Because the variables are set on the subprocess environment, any process the hook spawns SHALL also inherit them without additional configuration.

#### Scenario: Script invoked by hook can read WIP vars
- **WHEN** a hook is configured as `"./scripts/setup.sh"` and `setup.sh` reads `$WIP_REF_NAME`
- **THEN** `setup.sh` receives the correct ref name

### Requirement: WIP_* env vars do not override the existing environment
The injected `WIP_*` variables SHALL be added to the hook's environment alongside inherited system environment variables. Existing env vars (PATH, HOME, etc.) SHALL remain unchanged.

#### Scenario: PATH is preserved
- **WHEN** a hook runs via `sh -c`
- **THEN** the hook can invoke binaries on the system PATH without modification
