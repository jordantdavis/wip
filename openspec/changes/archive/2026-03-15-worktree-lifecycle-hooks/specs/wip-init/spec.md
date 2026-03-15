## MODIFIED Requirements

### Requirement: Init is idempotent
Running `wip init` multiple times from the same git repo root SHALL produce the same result as running it once. Subsequent runs SHALL NOT fail, overwrite state, or produce spurious output. This includes `.wip.yml`: if it already exists, it SHALL be left unchanged.

#### Scenario: Repeated invocation
- **WHEN** the user runs `wip init` twice from the same git repo root
- **THEN** both invocations exit with code 0, the second produces no output, and `.wip.yml` is not overwritten
