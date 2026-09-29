## Visible thread titles

### Choose the role and identity

Name a visible continuing thread for the work it actually owns:

| Role | Format | Use for |
| --- | --- | --- |
| Task | `Task[<identity>] :: <description> :: <YYYY-MM-DD>` | One issue or project task. |
| Owner | `Owner[<identity>] :: <description> :: <YYYY-MM-DD>` | A continuing owner that coordinates several tasks and checks their combined result. |
| Oracle | `Oracle[<identity>] :: <YYYY-MM-DD>` | An entity with deep project expertise that can speak on the project's behalf. |
| Auto | `Auto[<name>] :: <description> :: <YYYY-MM-DD>` | An automation the user requested or started. |

- Use `#<ticket>` only when the thread delivers that ticket or owns that epic. Otherwise use `$/<project-area>` with the user's established project name.
- Do not choose an identity from a passing ticket mention, branch, working directory, host, or the way the thread was created.
- One helper, one issue, or the word “owner” in a plan does not make a `Task[...]` an `Owner[...]`.
- A schedule or temporary helper does not make a thread an `Auto[...]`.
- An `Oracle[...]` explains the project's design, history, decisions, and current state. It need not coordinate delivery tasks. Omit the description and use the project's established `$/<project-area>` identity, as in `Oracle[$/leafnote] :: 2027-04-03`.

### Keep the title short and the date stable

- For Task, Owner, and Auto, use a short, specific description, usually three to eight words. Omit:
  - Host, branch, checkout, and home directory.
  - Duplicate dates and literal `<...>` placeholders.
  - Retired category prefixes and generated temporary titles.
- Keep the entire title within 59 characters, including spaces, punctuation, identity, and date. Count before applying it. Shorten the wording without losing the format, meaningful identity, or full `YYYY-MM-DD` date.
- Use the thread's actual creation date in the user's configured timezone. Preserve it when the ticket, role, description, or project changes.
- If the creation date is missing, read the thread metadata. Use older evidence only when it establishes the date, and say when the date is inferred. Do not substitute today's date.
- Preserve an established `Auto[...]` name and original date. Rename an automation or return an owner to `Task[...]` only when the user's direction or a verified role change requires it.

### Rename the existing thread

- When its identity or role changes, rename the same thread. Preserve its ID, original date, owner, host, goal, and useful history.
- For example, the same thread can become:
  1. `Task[$/leafnote] :: Replace search index :: 2027-04-03` for the project task.
  2. `Task[#APP-77] :: Replace search index :: 2027-04-03` when it delivers that ticket.
  3. `Owner[#EPIC-12] :: Lead search migration :: 2027-04-03` only when it becomes the continuing owner of that epic.
- Do not create a replacement thread to tidy its name.

To apply a title:

1. Identify the existing thread, its actual role, and its original creation date.
2. Choose a title that follows the format and length limit.
3. Use a tool that changes that same thread.
4. Read the thread again and verify the saved title. If a temporary race requires a retry, retry only on that thread.

If no title tool is available, report the thread ID, desired title, and missing tool. Do not claim that the rename succeeded.

### Keep naming separate from other actions

- Decide whether to create a thread, choose a machine, move work, or assign an owner separately, using the available runtime tools and the user's actual permission.
- Do not create, archive, hand off, move, or choose a host for a thread just to apply a title.
- A temporary helper without a visible thread does not need a title in this format.
