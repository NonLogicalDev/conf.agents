## Work delegation

- Keep owner tasks responsive to the operator, stakeholders, steering, and
  changes in direction. Prefer coordination and project management, but
  let an owner do the work directly when one worker is the simplest fit.
- Use [$Tasker_Plunger]({{%_skills_%}}/Tasker_Plunger/SKILL.md) when one owner coordinates several subagents and needs a durable project plan, shared owner context, or a clear handoff.
- Mark every task being worked as `active` and include its assigned subagent's verified full `/root/...` path. Reserve unmarked `pending` labels for work that has not started.
- Mark every running task `in_progress` when the Steps tool supports several active entries. Keep all active tasks visible and annotated with their workers.
- Use effort estimates as guidance, not hard limits. Choose the worker
  type that best fits the task:
  - **Subagents:** Prefer these for independent work expected to take more
    than one minute. Give each helper a clear scope and distinct files.
  - **Local task threads:** Consider these for self-contained work likely
    to take more than 20 minutes and benefit from a persistent local owner.
  - **Remote task threads:** Consider these for similarly substantial work
    that benefits from a verified, available, user-allowed remote machine.
- Reuse an existing matching task owner. Create a new user-visible task
  thread only when the operator explicitly requests one. A local task
  running a remote shell is not a remote task.
- Keep the original owner accountable for assignments, progress,
  stakeholder communication, result verification, and final integration.
  Use the available runtime delegation tools when separate workers help.
- Consult the operator when the best owner, worker, or execution location
  is genuinely unclear.
