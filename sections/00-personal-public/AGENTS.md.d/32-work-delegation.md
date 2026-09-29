## Work delegation

- Keep the owner responsive to the operator, stakeholders, and changes in direction. Prefer coordination and project management; the owner may do the work directly when one worker is the simplest fit.
- Use [$@:Tasker_Stream](skill://@:Tasker_Stream) when one owner coordinates several subagents and needs a durable project plan, shared owner context, or a clear handoff.
- Keep the current assignments visible:
  - Mark work currently being done as `active` and include the assigned subagent's verified full `/root/...` path.
  - Use `pending` without an active marker only for work that has not started.
  - When the Steps tool supports several active entries, mark every running task `in_progress` and name its worker.
- Choose a worker suited to the task. Treat these time estimates as guidance, not limits:
  - **Subagents:** Prefer these for independent work expected to take more than one minute. Give each a clear scope and distinct files.
  - **Local task threads:** Consider these for work that can stand on its own, is likely to take more than 20 minutes, and would benefit from a persistent local owner.
  - **Remote task threads:** Consider these for similarly substantial work when a remote machine would help. Verify that the machine is available and the user permits its use.
- Reuse an existing owner when the work fits. Create a new visible task thread only when the operator explicitly requests it. Running a remote shell from a local task does not make that task remote.
- The original owner remains responsible for assignments, progress, communication with stakeholders, checking results, and combining the work. Use the available runtime tools to delegate when separate workers help.
- Ask the operator when the choice of owner, worker, or machine is genuinely unclear.
