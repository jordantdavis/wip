## ADDED Requirements

### Requirement: Hooks execute via sh -c
Every hook command string in a `wip` hook list SHALL be executed as `sh -c "<hook>"` rather than by splitting the string and invoking the binary directly. This applies to all current and future hook types.

#### Scenario: Simple command runs correctly
- **WHEN** a hook is configured as `"npm install"`
- **THEN** the CLI executes it as `sh -c "npm install"` and the command runs successfully

#### Scenario: Compound command with && runs correctly
- **WHEN** a hook is configured as `"npm install && npm run build"`
- **THEN** both commands execute in sequence and the hook succeeds only if both exit with code 0

#### Scenario: Hook with pipe runs correctly
- **WHEN** a hook is configured as `"cat package.json | jq .version"`
- **THEN** the shell pipe executes correctly

#### Scenario: Hook with output redirect runs correctly
- **WHEN** a hook is configured as `"npm install > install.log 2>&1"`
- **THEN** output is redirected to the file without error

#### Scenario: Shell variable expansion works in hooks
- **WHEN** a hook is configured as `"echo $WIP_REF_NAME"`
- **THEN** the shell expands the variable and the ref name is printed
