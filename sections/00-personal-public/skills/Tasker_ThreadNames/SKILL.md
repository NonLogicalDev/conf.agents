---
name: Tasker_ThreadNames
description: Choose, update, and check a clear Task[...], Owner[...], or Auto[...] title for a task the user can see. Keep its existing owner, host, and original start date.
---

# Tasker Thread Names

Use this skill to choose and verify the name of a thread the user can see. It does
not create a thread, choose a machine, move work, or assign an owner. Make
those decisions separately using the available runtime tools and the user's
actual permission.

## Choose The Right Name

Use the role the thread actually owns:

| Role | Title |
| --- | --- |
| One issue or project task | `Task[<identity>] :: <description> :: <YYYY-MM-DD>` |
| A project owner responsible for continuing work | `Owner[<identity>] :: <description> :: <YYYY-MM-DD>` |
| An automation the user requested or started | `Auto[<name>] :: <description> :: <YYYY-MM-DD>` |

Use `#<ticket>` only when the thread actually delivers that ticket or owns
that epic. Otherwise, use `$/<project-area>`. A passing mention, a branch, a
working directory, or a host does not set a thread's identity.

Use `Owner[...]` only for a thread responsible for several tasks and their
combined result. Calling one helper does not turn a `Task[...]` into an
owner. Use `Auto[...]` only for an automation the user has actually started
or requested. A temporary helper does not need a visible thread title.

Read [references/title-rules.md](references/title-rules.md) for examples,
project names, title changes, original dates, and missing rename tools.

## Keep The Original Thread And Date

Use the thread's real creation date in the user's configured timezone. Keep
that date when its role, ticket, description, or project changes. Read the
actual creation record if the original date is missing. Do not guess today.

When a role changes, update the same thread. Preserve its ID, original date,
owner, host, goal, and history. Do not create a replacement or switch hosts.

## Apply And Check The Title

1. Identify the existing thread, its actual role, and its original date.
2. Choose the exact title and a short description of the actual work.
3. Use an actual tool that changes the title of that same thread.
4. Read the same thread again and check the exact saved title.
5. If the tool is unavailable, report the thread ID, requested title, and
   missing tool. Do not claim that the title changed.

Do not create, archive, hand off, or move a thread to apply a title.

## Tests

Read [tests/README.md](tests/README.md). Run the relevant tests with fresh,
independent agents when available.
