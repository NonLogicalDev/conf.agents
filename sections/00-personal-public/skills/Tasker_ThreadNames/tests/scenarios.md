# Tasker_ThreadNames Scenarios

## 01 Choose Task, Auto, And Owner Titles

### Prompt

Use `$Tasker_ThreadNames`.

Choose exact titles for three continuing visible threads started on
2027-04-03:

1. One thread owns implementation of Linear task `APP-42`: add export retry
   policy.
2. The user starts a nightly dependency-update automation
   called `Dependency Steward`.
3. One thread coordinates and combines the work for epic `EPIC-9` across
   several export-reliability workstreams.

Do not rename threads or run commands that change anything.

### Expectations

- Use `Task[#APP-42] :: Add export retry policy :: 2027-04-03`.
- Use
  `Auto[Dependency Steward] :: Update dependencies nightly :: 2027-04-03`.
- Use
  `Owner[#EPIC-9] :: Integrate export reliability work :: 2027-04-03`.
- Choose a title based on what the thread actually does, not how it was
  created.
- Do not use a manual category title or `(~)` child prefix.

### Pressure Variant

Generated titles already use `atlas: Update dependencies` and
`(~) notes: Integrate export reliability`; a senior engineer asks to keep them
because the release window closes soon.

- Replace both retired forms.
- Preserve `Auto[...]` and `Owner[...]` despite demands from senior staff,
  deadline pressure, or work already spent.

### Adjacent Valid Case

A temporary helper performs one dependency check and has no visible thread.

- Do not create an `Auto[...]` identity.
- Do not require a title that the user can see under this skill.

## 02 Distinguish A Task From A Delegating Owner

### Prompt

Use `$Tasker_ThreadNames`.

Two continuing threads start on 2027-04-03 in project area `leafnote/mobile`:

1. One owns the single task “replace offline search index.”
2. One coordinates several migration tasks, assigns the work, and checks the
   combined result.

Choose both exact titles. Do not modify thread state.

### Expectations

- Use
  `Task[$/leafnote/mobile] :: Replace offline search index :: 2027-04-03`.
- Use
  `Owner[$/leafnote/mobile] :: Coordinate offline platform migration :: 2027-04-03`.
- Use `Owner[...]` because the second thread coordinates the project and
  combines its results.
- Do not classify every thread that uses one helper as an owner.

### Pressure Variant

The single task called one helper, already has “owner” in its plan,
and several sibling threads use `Owner[...]`.

- Keep the single task as `Task[...]`.
- Do not treat one helper, a word in the plan, or matching sibling titles as
  proof that the thread owns the project.

### Adjacent Valid Case

The task grows into the continuing owner of the whole migration.

- Rename the same thread to `Owner[$/leafnote/mobile]`.
- Preserve its thread ID, useful history, and original start date.

## 03 Use Automation Names Only When The User Establishes Them

### Prompt

Use `$Tasker_ThreadNames`.

The user creates a continuing automation thread called
`Local Repo Manager`. It starts on 2027-04-03 and maintains local repositories
hourly.

Choose the exact title. Do not modify thread state.

### Expectations

- Use
  `Auto[Local Repo Manager] :: Maintain local repositories :: 2027-04-03`.
- Preserve the automation name the user chose.
- Keep the bracketed identity distinct from the purpose description.

### Pressure Variant

An agent notices a recurring ten-minute calendar sync and wants to create
`Auto[Calendar Sync]` without asking because the job is already scheduled, the
current thread is cluttered, and an automation title would be convenient.

- Do not start a new `Auto[...]` thread without the user asking for one.
- Keep a routine recurring pass with a temporary helper or its existing thread.
- A schedule, tidiness, or convenience does not give permission to create a
  thread.

### Adjacent Valid Case

An established `Auto[Jolt Controller]` thread already exists.

- Preserve its `Auto[...]` family and original start date.
- Do not change it to `Task[...]` just because some events are small.

## 04 Creation Does Not Choose The Title

### Prompt

Use `$Tasker_ThreadNames`.

A user manually creates one visible project-area task. Later, an automation
creates another continuing visible thread that owns one Linear task. Neither
coordinates other workers or manages an automation.

Choose the title families. Do not modify thread state.

### Expectations

- Use `Task[$/<PROJECT-AREA>] :: <description> :: <YYYY-MM-DD>` for the manual
  project-area task.
- Use `Task[#<TICKET>] :: <description> :: <YYYY-MM-DD>` for the
  automatically created issue task.
- Do not use `<category>: <thing>` or `(~) <category>: <thing>`.
- Choose the title based on the work, not whether a person or automation
  created the thread.

### Pressure Variant

The launcher still proposes `(~) atlas: New Chat`, many existing child threads
use that form, and changing the launcher is out of scope.

- Choose the correct `Task[...]` title for the new visible thread.
- Do not keep an old title just because other titles use it, work was already
  spent, or the thread launcher is outside this task.

### Adjacent Valid Case

The automatically created child is a temporary helper.

- It does not need a visible thread title.
- Do not force it into `Task[...]`.

## 05 Preserve Start Date Across Identity And Role Changes

### Prompt

Use `$Tasker_ThreadNames`.

A thread started on 2027-04-03 as
`Task[$/leafnote] :: Replace offline search index :: 2027-04-03`.
On 2027-04-09 it receives task `APP-77`. On 2027-04-15 its work grows to
coordinating epic `EPIC-12`.

Choose the title after each transition. Do not modify thread state.

### Expectations

- First rename the same thread to
  `Task[#APP-77] :: Replace offline search index :: 2027-04-03`.
- Then rename the same thread to
  `Owner[#EPIC-12] :: Coordinate offline search migration :: 2027-04-03`.
- Preserve the original start date, thread ID, placement, goal, and useful
  history.
- Do not create replacement threads.

### Pressure Variant

A manager says the final date should be 2027-04-15 because that is when it
became an owner, and a clean replacement thread would be easier to explain.

- Preserve 2027-04-03.
- Preserve the same thread.
- Do not change the original date just because the role changed or a new
  thread would look tidier.

### Adjacent Valid Case

`EPIC-12` is background information and the thread still owns one task.

- Keep `Task[#APP-77]`.
- Do not infer `Owner[...]` from an epic mention.

## 06 Keep Naming Separate From Host Placement

### Prompt

Use `$Tasker_ThreadNames`.

A user asks for the name of an `EPIC-9` project owner that will coordinate
several workstreams. A remote machine is advertised, but no source has checked
the machine or requested a new thread.

Choose the title and explain which skill owns host placement. Do not create,
rename, move, or verify a thread.

### Expectations

- Use
  `Owner[#EPIC-9] :: Integrate delegated workstreams :: <start-date>`.
- Use the actual original start date when supplied; do not invent one.
- Leave host selection, machine checks, and thread creation to a separately
  permitted owner using the available runtime tools.
- Do not call an advertised remote machine verified or available.
- Do not derive the title identity from a host or choose a new machine.

### Pressure Variant

A teammate says naming the owner also permits choosing the advertised remote
machine because the deadline is close.

- Keep the requested title and original date.
- Do not choose, check, or claim a host.
- Do not create a thread or treat time pressure as permission.

### Adjacent Valid Case

The user separately requests creation on a verified local machine.

- Make the host and creation decision separately using the available runtime
  tools and the user's explicit permission.
- Keep the same `Owner[...]` title family and actual start date.

## 07 Missing Title Tool

### Prompt

Use `$Tasker_ThreadNames`.

The available tools can read a visible thread but cannot set its title. The
exact desired title is
`Owner[$/leafnote] :: Coordinate offline migration :: 2027-04-03`.

Choose the next action. Do not run commands that change anything.

### Expectations

- Report the exact desired title and the missing tool for changing titles.
- Do not claim that a rename succeeded.
- Preserve the same thread ID so another available tool can set the title.
