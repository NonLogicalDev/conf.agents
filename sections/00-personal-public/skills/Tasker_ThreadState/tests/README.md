# Tasker Thread State Behavioral Tests

Run each test with a fresh helper that has no prior context when the available tools allow it. Tell the helper to use `$Tasker_ThreadState` and give it the prompt, but not the expected answer.

Do not let tests change files unless a test names a temporary output file for that task. Save the helper's answer and compare it with the expected results afterward. A test passes only when it meets every expected result and does nothing it was told not to do.

When fixing a behavior, run the current instructions first. Run the same test again after the edit, including the pressure and nearby valid cases. Check that the agent follows the exact layout the user requests without guessing owners, links, times, blockers, or completed results.

Check that every running workstream is marked `active` with each assigned subagent's verified full `/root/...` path. Keep unstarted work `pending`, preserve real blockers, and never invent an unknown worker identity.
