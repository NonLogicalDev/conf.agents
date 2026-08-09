# Thread Title Rules

Use this reference when choosing or checking a visible thread title. The
naming skill does not select a host or create a thread.

## Task Titles

A task owns one issue or one project result:

```text
Task[#APP-42] :: Add export retry policy :: 2027-04-03
Task[$/leafnote/mobile] :: Replace offline search index :: 2027-04-03
```

Use the ticket only when the thread actually implements it. Research about a
ticket does not change the thread's identity. Use a project name already
established by the user when no ticket owns the work.

## Owner Titles

An owner coordinates work and checks the combined result:

```text
Owner[#EPIC-9] :: Integrate export reliability work :: 2027-04-03
Owner[$/leafnote/mobile] :: Coordinate offline migration :: 2027-04-03
```

Use an epic ticket only when the owner is actually responsible for that
epic. Otherwise, use the established project area. One helper or one issue
does not make a thread a project owner.

## Automation Titles

Use an automation title only when the user starts or requests that automation:

```text
Auto[Dependency Steward] :: Update dependencies nightly :: 2027-04-03
Auto[Local Repo Manager] :: Maintain local repositories :: 2027-04-03
```

Preserve the user's automation name. A scheduled command or a temporary
helper does not create a new automation thread.

## Keep Names Short And Stable

- Use a short, specific description, usually three to eight words.
- Use the user's established project name, including useful dots or hyphens.
- Do not include the host, branch, checkout, home directory, or date twice.
- Do not place literal `<...>` placeholders in the saved title.
- Do not use an old category prefix or a generated temporary title.

## Change The Same Thread

If a project task receives its own ticket, update the same task:

```text
Task[$/leafnote] :: Replace offline search index :: 2027-04-03
Task[#APP-77] :: Replace offline search index :: 2027-04-03
```

If that same thread becomes responsible for an epic, update its role without
changing its ID or creation date:

```text
Owner[#EPIC-12] :: Coordinate offline search migration :: 2027-04-03
```

Keep the existing owner, host, goal, useful history, and start date. If a
thread no longer coordinates a project, return it to `Task[...]` only when
that role change is real. Do not rename an existing `Auto[...]` without the
user's direction.

## Check The Date And Saved Title

Use the user's configured timezone for the original creation date. If the
date is missing, use thread metadata. Use older evidence only when it truly
establishes the creation date and explain that it is inferred.

Set the title on the existing thread ID. Read back that same ID and check the
exact saved title. Retry a temporary race only on that same thread. If no
title tool is available, report the ID, desired title, and missing tool.

When the user separately asks who owns work, where a task should run, or
whether to create a visible thread, use the available runtime tools and the
permission that request actually provides.
