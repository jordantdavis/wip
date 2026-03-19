## 1. Shared Validation and Path Logic

- [x] 1.1 Replace `validateWorktreeName` (regex) in `cmd/worktree.go` with `validateBranchName` that runs `git check-ref-format --branch <name>` and returns an error if it exits non-zero
- [x] 1.2 Add a `worktreePathSegment(branchName string) string` helper in `cmd/worktree.go` that replaces all `/` with `-`

## 2. worktree add

- [x] 2.1 Update `worktreeAdd` to call `validateBranchName` instead of `validateWorktreeName`
- [x] 2.2 Update `worktreeAdd` to derive `absWorktreePath` using `worktreePathSegment(worktree)` instead of the raw `worktree` argument

## 3. worktree remove

- [x] 3.1 Update `worktreeRemove` to call `validateBranchName` instead of `validateWorktreeName`
- [x] 3.2 Update `worktreeRemove` to derive `absWorktreePath` using `worktreePathSegment(worktree)` instead of the raw `worktree` argument

## 4. worktree list

- [x] 4.1 Update `worktreeList` to read the branch name for each worktree by running `git -C <abs-path> branch --show-current`
- [x] 4.2 Update `worktreeList` output to three columns: `<submodule>  <path-segment>  <branch>`

## 5. Tests

- [x] 5.1 Update any existing tests that assert on the old regex validation error message to match the new `git check-ref-format` error path
- [x] 5.2 Add tests for `validateBranchName`: valid simple name, valid slash name, invalid name
- [x] 5.3 Add tests for `worktreePathSegment`: no slashes, single slash, multiple slashes
- [x] 5.4 Add tests for `worktree list` output format (three-column)
