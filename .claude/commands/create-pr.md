Create a pull request for the current branch.

## Steps

1. Determine the default branch of the repository by running `gh repo view --json defaultBranchRef --jq '.defaultBranchRef.name'`.
2. Fetch the latest default branch from the remote: `git fetch origin <default-branch>`.
3. Run `git log origin/<default-branch>..HEAD --oneline` to see all commits on this branch.
4. Run `git diff origin/<default-branch>...HEAD --stat` to see which files changed.
5. Run `git diff origin/<default-branch>...HEAD` to read the actual changes.
6. Analyze the changes to determine:
   - **Type**: `feat`, `fix`, `refactor`, `perf`, `deprecate`, `docs`, `test`, or `chore`
   - **Component**: infer from the area of the codebase that changed. Use a short, lowercase name (e.g. `transfers`, `accounts`, `cards`, `client`). If the change spans the whole project, omit the component.
   - **Breaking change**: whether this could break existing consumers or dependents. If yes, add `!` after the component.
7. Generate the PR title in this format:
   ```
   <type>(<component>)[!]: <description>
   ```
   Rules:
   - Description must be imperative mood, lowercase, concise
   - Focus on what the PR does, not how
   - `chore`, `docs`, `test` do not require a component
   - `!` is required for breaking changes
   - Do not include Linear ticket references (e.g., CAR-*, LEDG-*) in the title
8. Determine which GitHub labels to apply:
   - `feat` -> `enhancement`
   - `fix` -> `bug`
   - `docs` -> `documentation`
   - `refactor`, `perf`, `deprecate`, `test`, `chore` -> no default GH label (skip `--label` for these)
   - If the PR has `!` (breaking change), add `breaking-change` label. If the label doesn't exist on the repo, skip it — don't fail.
9. Generate the PR description using this template:

   ```
   ## What

   <!-- 1-3 bullet points: what changed -->
   - ...

   ## Why

   <!-- Brief explanation of motivation or context -->

   ## Notes

   <!-- Optional: anything reviewers should know, migration steps, breaking change details. Delete if not needed. -->
   ```

   Fill in the template:
   - Prefix the entire description with 🤖 as the first character.
   - Write the description in ASD-STE100 Simplified Technical English.
   - "What" — bulleted list of concrete changes
   - "Why" — motivation: what problem this solves or what it enables
   - "Notes" — only include if there are breaking changes, migration steps, or important reviewer context. Delete the section entirely otherwise.
   - End with a signature line naming the agent, e.g. `— Claude`.

10. Before creating the PR, show the user the generated title, labels, and description. Ask for confirmation.
11. Create the PR using `gh pr create --title "TITLE" --body "BODY"` with the generated title, description, and labels.
    - Use `--title` and `--body` (or `--body-file`) so the command is not interactive.
    - Use `--label` flag for each label (e.g., `--label enhancement`)
    - Push the branch first if needed (`git push -u origin HEAD`)
    - If a label doesn't exist on the repo, retry without that label.
12. Return the PR URL.

## Label mapping

| Type | GitHub Label |
|---|---|
| `feat` | `enhancement` |
| `fix` | `bug` |
| `docs` | `documentation` |
| Any type with `!` | also add `breaking-change` (if it exists on the repo) |

## Examples

Good PR titles:
```
feat(transfers): add cancellation client methods
fix(accounts): return 404 instead of 500 for missing resources
refactor(client): extract bearer-token header setup
fix(transfers)!: require account-prefixed transfer paths
docs: update installation snippet
chore: upgrade Go to 1.26
test(mhooks): add signature mismatch cases
deprecate(accounts): mark UpdateAccount for removal
```
