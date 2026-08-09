## Code projects and local checkouts

- Treat the session's starting directory as coordination context, not proof of which project is being changed. A session may start in `<agent-brain>/` while the actual work belongs to a different checkout.
- Before any project work, resolve the actual project, checkout or worktree,
  and target files. Find and read every instruction file that applies to
  that work:
  - Start with the global instructions.
  - Read the project or worktree root's `AGENTS.md` and each applicable
    `AGENTS.md` between that root and the target directory.
  - Read `AGENT.md` too when the project uses that filename.
  - Check again after changing projects, worktrees, or target directories.
  - For a new checkout, read its instructions immediately after cloning
    and before continuing.
- Before inspecting, fetching, or cloning a repository, check whether its
  existing local checkout can be resolved from the request and applicable
  global or repository instructions.
- Use `<projects-home>` as the configured project root; `~/projects/` is a suggested default:
  - Keep reusable remote checkouts under
    `<projects-home>/remote/<host>/<namespace...>/<repo>/`.
  - Keep local-only projects under
    `<projects-home>/local/<phase>/<type>.<name>/`.
    - Use `active` for projects intended for continued work.
    - Use `archived` for retired projects and experiments.
    - Common types include `app`, `cli`, and `lib`. Extended types are
      supported.
- Derive a remote checkout from its full repository URL. Preserve the host
  and every namespace segment. Treat HTTPS and SSH URLs for the same
  repository as the same checkout.
  - `<host>` is the Git server, such as `github.com`.
  - `<namespace...>` is the repository owner, organization, or complete
    sequence of nested groups.
  - `<repo>` is the repository name without a trailing `.git`.
  - For example, both
    `https://github.com/NonLogicalDev/gymnasium` and
    `git@github.com:NonLogicalDev/gymnasium.git` map to
    `<projects-home>/remote/github.com/NonLogicalDev/gymnasium/`.
  - Apply the same mapping to other Git hosts and preserve all nested
    namespace segments.
- Check a user-identified existing checkout first. Apply a more specific
  configured host, organization, project, or destination rule when present.
- Check the exact resolved path. Do not enumerate project roots or unrelated
  checkouts unless the user explicitly requests a bounded inventory.
- Reuse an existing checkout. If the task requires current remote state,
  verify its freshness instead of assuming local files are current.
- If no checkout exists, use the resolved reusable project path. Use a
  permitted temporary location for one-off inspection.
