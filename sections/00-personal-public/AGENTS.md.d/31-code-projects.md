## Code projects and local checkouts

### Find the checkout before working

- Identify the actual project, checkout or worktree, and files involved in the task. The starting directory does not establish which project owns the work; a session in `<agent-brain>/` may need a different checkout.
- Before inspecting, fetching, or cloning a repository, use the request and applicable global or repository instructions to look for an existing checkout.
- Check an existing checkout named by the user first. Follow any more specific configured rule for the host, organization, project, or destination.
- Check the resolved path directly. Do not list project roots or unrelated checkouts unless the user explicitly requests an inventory with a defined scope.
- Reuse an existing checkout. When the task needs current remote state, check freshness instead of assuming local files are current.
- If no checkout exists, use the reusable path described below. A permitted temporary location may be used for a single inspection.

### Read the applicable instructions

Before project work:

1. Read the global instructions.
2. Read `AGENTS.md` at the project or worktree root and in each applicable directory between that root and the files involved in the task.
3. Read `AGENT.md` too when the project uses that filename.

Read the applicable instructions again when changing projects, worktrees, or directories. For a new checkout, read them immediately after cloning and before continuing.

### Use the configured project layout

- Use `<projects-home>` as the configured project root. `~/projects/` is a suggested default.
  - Keep reusable remote checkouts under `<projects-home>/remote/<host>/<namespace...>/<repo>/`.
  - Keep projects that exist only locally under `<projects-home>/local/<phase>/<type>.<name>/`.
    - Use `active` for projects intended for continued work.
    - Use `archived` for retired projects and experiments.
    - Common types include `app`, `cli`, and `lib`; other types are supported.
- Derive the remote checkout path from the full repository URL:
  - `<host>` is the Git server, such as `github.com`.
  - `<namespace...>` is the repository owner, organization, or full sequence of nested groups.
  - `<repo>` is the repository name without a trailing `.git`.
- Preserve the host and every namespace segment for every Git host. HTTPS and SSH URLs for the same repository must resolve to the same checkout.
- For example, both `https://github.com/NonLogicalDev/gymnasium` and `git@github.com:NonLogicalDev/gymnasium.git` map to `<projects-home>/remote/github.com/NonLogicalDev/gymnasium/`.
