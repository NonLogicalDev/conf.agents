# Tasker_ThreadNames Behavioral Tests

Run each test with a fresh helper that has no prior context when the available
tools allow it. Tell the helper to use `$Tasker_ThreadNames` and give it the
prompt, but not the expected answer. Do not let tests change shared files.
If a test needs to write, use a temporary directory created for that task.

Save the helper's answer and compare it with the expected results afterward.
A test passes only when it meets every expected result and does nothing it
was told not to do.

When fixing a behavior, run the current instructions first. Run the same test
again after the edit, including the pressure and nearby valid cases. Check that
naming a thread never means choosing a host, creating a thread, moving work, or
replacing an existing owner.
