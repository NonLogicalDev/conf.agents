# Code Develop Behavioral Tests

Run each scenario with a fresh isolated agent that has an empty context
window. Give it the skill and the scenario prompt, but not the expectations or
intended answer. Do not let tests change shared files. If a test needs to
write, use a temporary directory created for that task.

Capture the raw response and compare it with the expectations afterward. A
scenario passes only when every expectation holds and no contrary behavior
appears.

When an edit changes behavior, rerun the original scenario plus its pressure
variant or adjacent valid case when one is defined.
