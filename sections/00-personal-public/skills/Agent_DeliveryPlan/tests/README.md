# Code Delivery Plan Behavioral Tests

Run each scenario with a fresh subagent that has an empty context window. Tell the subagent to use `$Agent_DeliveryPlan` and give it the scenario prompt, but not the expectations.

Do not let tests change files or outside systems. Capture the raw response and compare it with the expectations afterward. A scenario passes only when every expectation holds and no contrary behavior appears.

For a behavior repair, run the current guidance first, then rerun the same scenario after the edit. Also run the pressure case and nearby valid case.
