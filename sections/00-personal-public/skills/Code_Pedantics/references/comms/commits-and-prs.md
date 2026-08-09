# Commits And Pull Requests

## Choose The Message Format First

Assume squash unless the repository or user explicitly states otherwise.

### Squash Workflow With A Single Local Commit

When the branch is squashed into one local commit, use the same reviewer prose for the commit and the pull request. Keep the commit subject identical to the PR title and the commit body identical to the canonical reviewer prose in the PR description. Append any required public-agent signature only to the PR description, after that prose. Keep the commit and PR title unsigned. Put other metadata required only by the PR in its designated fields.

### Squash Workflow With A Branch Series

Prefer one branch commit when the workflow allows it. When a squash workflow keeps several commits, give each new commit the current message for the entire branch diff. Keep the newest commit's subject and body identical to the PR title and canonical reviewer prose. Keep the required public-agent signature and other PR-only metadata out of the commit. Older commits may keep the messages that described the branch when they were created. Do not require them to match the current PR or describe only their individual changes.

### Non-Squash Workflow

When the user or repository explicitly chooses a non-squash workflow, write each commit message for that commit's own diff. Write the PR as a summary of the complete series. Follow the project's format or the default below. Keep its purpose, behavior, checks, limitations, and risks consistent with the commits. The PR does not need to match any individual commit.

### Metadata Required Only By The Pull Request

When a required PR template prevents the complete descriptions from matching, put the canonical squash message in the corresponding PR title and body fields. Put additional metadata in the fields or sections meant for it, and keep any required public-agent signature as the final footer of the PR description. In a non-squash workflow, keep PR-only metadata and the signature separate from the summary of the series and the individual commit messages. Identify the mapping before publication.

### Keep Conversation State Out Of Reviewer Prose

A pull request description is a lasting reviewer document. Keep its reviewer prose free of private task links, drafting history, and other details about the conversation that produced it. A signature required by the applicable public-agent signature contract belongs after that prose, including its required source-thread link.

Apply the required signature to every communication that the contract covers unless the user explicitly asks not to sign that communication. This requirement overrides conflicting writing, formatting, template, and commit-parity rules. Keep signatures out of PR titles and Git commit messages.

### Keep A Published Pull Request Current

After amending, rebasing, changing scope, addressing review, or recording a new check, read the final diff again and update the chosen message. If the PR is already published, update and verify its metadata after every push:

- for a squash with one commit, update the canonical PR fields from the commit message;
- for a squash branch with several commits, use the newest commit's message for the entire branch;
- for non-squash, update the PR summary for the complete series; and
- when a required template uses separate metadata, update those fields too.

After the permitted publication workflow pushes, read the forge state again. Verify the matching reviewer prose, required signature, and any required template mapping. Do not claim that the push, review handoff, or publication is complete while the published PR still shows stale metadata. If you lack permission or a supported way to update the metadata, stop and report the blocker. Do not silently leave the mismatch behind.

These instructions do not give you permission to commit, publish, or write to a forge.

See [examples/commits-and-prs.md](examples/commits-and-prs.md) for examples of a squash with one commit, a squash with several commits, a non-squash workflow, and small changes.

## Scope

Use this reference for intentional commit messages, squash messages, pull request titles, pull request descriptions, and summaries written for reviewers. Use a stricter repository convention or required template when one exists, while preserving any required public-agent signature.

Use this reference to improve wording. It does not give permission to commit, publish, request review, send messages, or merge.

When the user asks for the commit message only, return only the commit message without a preamble or explanation. That request still means a complete subject and body. Omit the body only when the user explicitly asks for a subject or title only.

## Project Convention Comes First

Before choosing a title format, body shape, headings, tense, or trailers, read the repository's instructions and inspect strong recent non-merge commits for the same area. Prefer commits that touch the same files or subsystem over a generic style guide.

Project conventions may define an area prefix, capitalization, subject mood, line length, narrative paragraphs, required trailers, pull request template, or named sections. Follow those conventions. The content contract below still applies, but summary and verification describe information the reviewer needs; use the project's headings when it already follows a different form.

When the project has no clear convention, use the fallback guidance here. Do not impose a generic semantic title or Markdown template on a project whose history uses a different reviewable form.

Do not splice together conflicting generic commit templates. Use this order: repository instructions and nearby history, then the user's explicit format, then this default house format. Carry over compatible content rules from another generic guide, but do not combine its headings or bullet rules with this format.

## Default House Format

When the project allows the author to choose the form, use:

```markdown
<semantic subject>

<one or two short unheaded paragraphs with the intent and crisp overview>

## What Changes

- <guide the reviewer through the important behavior, responsibility, or
  contract changes>

## Verification

- `<repeatable command or check>`
  - <what it proved>

<optional trailer required or preferred by the user or project>
```

Keep the opening one or two paragraphs extremely crisp. `## What Changes` should use clear bullets and explain the overall change without reciting exact code edits. `## Verification` should describe validation that actually ran, with repeatable commands or checks, the behavior they proved, and why that evidence supports confidence in the work. Avoid transient run links or bare commit hashes as verification evidence; the durable description of the check should carry the claim.

When verification evidence is missing, use one clean status bullet:

- `Not run: no verification was performed.` when no check ran;
- `Unknown: verification status is unavailable.` when the status is unknown; or
- `Not applicable: <why no meaningful check exists>.` when there is genuinely nothing to verify.

Do not append prompt-supplied check names such as lint, formatting, links, or tests unless their status is materially relevant to review.

In the default house format, include `## Verification` by default. Include `## What Changes` for any change whose overall diff needs reviewer guidance. When the diff is small and its meaning is obvious, omit `## What Changes` if the opening summary already carries the complete reviewable story. Do not use the exception to hide several responsibilities or a change that needs explanation.

`## Verification` is the default. Omit it only as a rare edge case after establishing that there is genuinely no meaningful behavior, build, rendering, link, command, or other claim to verify. If that judgment is not obvious and defensible, keep the section and use a clean `Not applicable` status instead.

## Give Reviewers The Context They Need

Every intentional commit or pull request needs enough context for an engineer to review it without hidden conversation history.

Write every commit and pull request so it makes sense on its own. Assume the reader has not seen the task conversation, review thread, issue discussion, or earlier draft. Include the context needed to understand the purpose, overall change, and proof in the artifact itself. Durable links may supplement the explanation, but they cannot be the only home for essential context.

Do not expose internal task plans, plan numbers, task conversation, unrelated source-thread links, private task IDs, chat state, coordinator instructions, hidden sequencing rationale, agent reasoning, or details about the user, agent, prompt, or drafting process as the source of the reviewer prose. Keep the source-thread link required by an applicable public-agent signature in that signature only. Explain the actual behavior, compatibility requirement, protocol, risk, or evidence the reviewer can verify. If the artifact is already published and the user permits a correction, do not leave private material visible just to avoid amending the message, updating the stack, or changing the head SHA.

Strive to make the commit maximally understandable to an intelligent reader with no prior relationship to the project. Explain unfamiliar project terms, the people or systems involved, their responsibilities, and why the change matters. Do not assume local shorthand, team folklore, or repository archaeology. Include the context the reader needs and leave out the rest.

Aim for a reviewer to understand the purpose, overall change, and verification in two to five minutes. A justifiably wide-reaching change may take up to ten minutes. If the necessary context cannot fit into one or two paragraphs or bullet groups per section, the change is probably too broad and should be split.

Include this content by default:

1. **Summary:** the intent behind the change, the problem or opportunity, and the overall behavior or architecture it changes.
2. **Verification:** concrete evidence that the change behaves as intended and does not break relevant behavior, written so another engineer can repeat the checks; when evidence is unavailable, an honest not-run, unknown, or not-applicable status with a reason instead of an invented claim. Omit this content only for a justified, genuinely not-applicable edge case.

In the default commit form, keep essential context inside the opening summary and, when needed, `## What Changes` instead of adding more headings. If it does not fit, the commit probably needs to split. A project template or pull request description may require a compact **Context**, **Root cause**, **Design**, **Contract**, **Migration**, or **Risk** section.

## Choose What Each Commit Includes

Before writing or approving a commit message, identify the exact intended change.

For the normal squash workflow, describe the final squash diff, even when several branch commits prepare it. In a non-squash workflow, make each commit's contents match its own message.

- Split unrelated purposes into separate commits.
- Exclude unrelated user changes.
- Keep recovery checkpoints separate from intentional public history.
- For non-squash commits and a squash with one commit, include only the behavior described in the message. When several commits prepare a squash, write each new message for the complete branch diff.
- If a message is too long, consider splitting the change into smaller commits.

## Gather The Facts

When the checkout is available, inspect the exact intended diff before drafting. Use the staged diff for a staged commit. Use the working tree diff only when those changes are intended for the commit. Check status and diff scope for unrelated changes, then use repository instructions, public interfaces, issue context, and actual verification output to fill the message.

Do not infer rationale or verification from filenames, helper names, changed test files, or the task conversation when the underlying evidence is available. If evidence is missing, say so rather than inventing it.

Treat prompt-supplied file names, folder names, helper names, fixture names, and check names as candidate details, not required commit content. Include them when they clarify the overview, explain an important component or behavior, or identify an actual check. Do not repeat them merely because the prompt lists them.

For a documentation change, a README, guide, folder, or example name can be useful orientation, but it does not replace the reader task or outcome. Pair a helpful artifact reference with the behavior it explains, and prefer a succinct overview over a list of every file.

## Titles

Follow the repository's required title format and inspect nearby non-merge commits when the format is implicit. Many projects use an area or filename prefix; others use semantic types, issue prefixes, or plain imperative subjects. Match the established form.

When the repository has no stronger convention and accepts semantic titles, use:

```text
<type>(<scope>): <imperative subject>
```

Use the narrowest accurate type and exact subsystem, package, module, command, or service as the scope. Avoid vague scopes such as `misc`, `cleanup`, or `updates`. In the fallback form, start the subject with an imperative verb, name the behavior or artifact, keep it at or below 72 characters when practical, prefer 50 characters or fewer, and omit a final period.

Follow the project's body-wrapping convention. When no convention exists, use semantic line breaks and keep prose lines at or below 72 characters when practical. Do not damage an unbreakable command, URL, or identifier merely to meet the fallback width.

This fallback applies to actual Git commit messages. Do not use it to hard-wrap pull request descriptions, review requests, documentation, task replies, or other Markdown prose.

## Summary Content

Include a useful summary even when the title is good. In the default house format, write it as one or two short unheaded paragraphs. Use the project's native heading or narrative form when the project has a stronger convention.

Start with what changes, why the change is needed, and what outcome it produces. For a technical change, name the current behavior or failure, explain the problem, and describe the chosen approach and its result. Explain which responsibilities, contracts, or behavior changed. Name the behavior and design instead of listing files, functions, comments, and tests.

Name exact public commands, flags, config keys, API fields, file formats, and visible errors when they are part of the changed contract. A reader new to the project should not have to infer which interface the message means.

A useful default summary has one or two short paragraphs. For a larger change, put guidance about the overall change in `## What Changes` rather than stretching the opening paragraphs into a transcript.

Good summary content:

- the problem experienced by the user or operator;
- the old behavior and why it was wrong or insufficient;
- the new behavior and the intended outcome;
- the chosen approach and why it fits better than obvious alternatives;
- who owns the behavior, which versions remain compatible, or which safety rule applies; and
- deliberate non-goals that keep the change reviewable.

Do not replay discovery history, rejected drafts, negotiation, or implementation chronology unless a lasting decision depends on it.

## Explain What Changes

In the default house format, use `## What Changes` as a short guide to the important behavior, responsibilities, contracts, compatibility, and migration. Write each point as a short bullet that tells a reviewer what to inspect and why it matters.

Do not turn the section into an inventory of files, helper functions, fixtures, comments, or test files. Mention a file, folder, or implementation detail when it helps locate a meaningful contract, ownership decision, compatibility requirement, migration, or surprising design choice.

Even when the user explicitly asks for file, helper, fixture, or routine-process inventory, keep only details that help review the change. Translate the rest into the behavior, responsibility, or reader task the change affects.

## Verification Content

Verification content is required by default. State what actually ran and what the result proved. In the default house format, use `## Verification`; otherwise include the same concrete evidence in the project's native body form.

Use Verification to justify confidence in the quality of the work, not merely to dump commands. Explain which changed behavior, risk, contract, or compatibility requirement each check verifies. State what the checks do not prove. Keep claims proportional to the evidence.

Include commands when they make the evidence repeatable or strengthen the case. Pair each useful command with the result and the behavior it covers; avoid only an unconnected transcript that leaves the reviewer to infer why it matters.

For a substantial or risky change, ask an independent engineer or agent who has not seen the work to review only the commit message, the intended diff or checkout, and the stated verification. Do not provide task thread context or the expected conclusion. The cold reviewer should be able to reconstruct the intent, inspect the change, and reach a conclusion supported by evidence of whether the change is safe, achieves its goal, and avoids obvious regressions. If the reviewer cannot, improve the message, run more relevant checks, or split the change.

Never invent verification. When no check ran or the result is unknown, keep the section and state that limitation plainly, including why verification is missing when known.

When evidence is unavailable, emit one clean status bullet and do not repeat prompt-supplied check names unless their status is materially relevant to the review.

Use “Not applicable: <reason>” when there is genuinely no meaningful behavior, build, rendering, link, or other claim to verify. Never leave “Not applicable” bare. Use “Not run” when a relevant check exists but was not performed. Do not label verification not applicable merely to avoid naming a skipped check.

Before omitting `## Verification`, articulate why no meaningful verification exists. In a review or explanation, state that justification. When the user asks for the commit message only, keep the justification out of the returned artifact and omit the section only when the summary makes the edge case clear.

Do not omit `## Verification` because evidence exists in a task thread, review comment, CI page, or other external artifact. Keep the commit understandable on its own; omission is allowed only when there is no meaningful verification to report.

Do not enumerate requested checks or fill the section with speculative statuses for lint, formatting, links, tests, or manual checks that were not material to the change.

Prefer “Unknown: verification status is unavailable” over “Unknown: lint and other verification status are unavailable.” The second form repeats an irrelevant requested check instead of telling the reviewer anything useful.

Prefer repeatable evidence:

- exact focused test, formatter, linter, type-check, build, or integration commands;
- the relevant environment, fixture, or target when it changes interpretation;
- a concise result such as passed counts or verified behavior;
- behavior proven through the actual entrypoint or complete workflow when the change depends on it; and
- honest limits, skipped checks, or blockers.

Changed test files are not verification evidence. A broad green status is not a substitute for naming the checks that cover the changed behavior. Do not rely on links to a single run or bare commit hashes. Repeatable commands, inputs, environment, and observed behavior should carry the claim.

Do not use routine formatter, linter, or generic “all checks pass” claims as a substitute for evidence about the changed behavior. Include routine hygiene only when it materially validates the change or the project requires it in the message.

For a large change, use one compact bullet group. Do not bury the reviewer in raw logs or every command from the development process.

## Remove Empty Meta And Hype

Remove framing only when it adds no meaning. Phrases such as “this change” are fine when they orient the reviewer or distinguish the diff; prefer a direct domain noun when it is clearer. Remove generic praise, inflated adjectives, and claims that are stronger than the evidence. Do not make a message sound more impressive at the cost of accuracy or reviewer time.

Prefer concrete verbs and exact domain nouns. A commit message should explain the behavior and proof, not advertise the patch or narrate that files, helpers, fixtures, comments, formatting, or tests changed.

## Pull Request Descriptions

Begin with what the pull request adds, fixes, removes, or changes. State the change and its effect in the first sentences. Put useful information about the service, product, or existing system in `## Background` after the summary. Include background only when it helps the reviewer understand the change. When the project uses `## Context` or another equivalent heading, follow its template instead.

A PR body must faithfully state the diff's intent, not only the resulting state. When the diff corrects or reverts work already on the default branch, say that plainly in the opening summary and name the earlier change when it is verified and helps review. State what was wrong or obsolete and what the diff restores or replaces. Do not make a correction sound like a clean new design or neutral cleanup. Do not invent a correction or revert relationship when evidence does not establish one.

A PR body should explain the change, not narrate its own status or the conversation that produced it. Leave out statements that the PR is a draft, ready for review, written as requested, or based on a task, thread, prompt, plan, operator instruction, agent decision, or earlier draft. Referring to the PR diff or saying `this change` is fine when it helps reviewers understand behavior, evidence, dependencies, or risks. Keep required repository template headings and fields.

Follow the rule for matching commit and PR messages at the top of this reference. For a squash workflow, keep the canonical reviewer prose in the PR description identical to the commit body unless a required template defines a mapping. Append any required public-agent signature only to the PR description. If the branch has several commits, match the newest commit's message for the entire branch. For a non-squash workflow, summarize the entire series and keep that summary consistent with the individual commits. Preserve the canonical content inside required template fields.

Write for a reviewer who cannot see the task conversation. Describe behavior and review decisions. Refer to the diff when it helps orient review, but do not let a file list replace the explanation. Keep each section to one or two paragraphs or compact bullet groups.

When a stack or dependency matters, state the relationship and the review order without narrating unrelated stack history. A public issue or parent pull request may appear when it explains a visible dependency or contract. Keep the link beside the fact it supports; do not use it instead of the explanation. When a change has a real limitation or unexecuted verification, state it plainly instead of hiding it behind a confidence claim.

## Check Before Review

Before handoff, verify:

- the title and body describe exactly the intended change;
- the intended staged or working tree diff was inspected when available;
- the commit stands on its own without task-thread or review-chat context;
- the commit and reviewer prose contain no private plan names, task numbers, conversation references, unrelated source-thread links, coordinator instructions, hidden sequencing rationale, agent reasoning, or authoring details about the user, agent, prompt, or draft;
- the pull request description includes its required public-agent signature and source-thread link unless the user explicitly exempted that communication;
- the pull request title and reviewer prose contain no agent signature, unrelated task link, narration about the PR's own draft or review status, or context about the conversation that produced it;
- the opening summary states when the diff corrects or reverts verified work already on the default branch instead of presenting the correction as neutral cleanup;
- a reader new to the project can understand the purpose and relevant terms;
- the summary begins with the change and its effect; useful service context appears in optional background;
- changed public interfaces use their exact names;
- verification content names real, repeatable evidence;
- missing or inapplicable verification is stated plainly rather than invented;
- omitted verification has a justified, genuinely not-applicable reason;
- the visible format follows the repository's instructions and nearby non-merge history;
- the commit or pull request excludes unrelated work;
- a squash branch with one commit and the canonical reviewer prose in its PR description are identical;
- a squash branch with several commits keeps its canonical PR title and reviewer prose identical to the newest commit's message for the entire branch;
- the commit subject, commit body, and PR title contain no agent signature or source-thread link;
- a non-squash PR consistently summarizes commits that each describe their own diff;
- the squash-path PR message is treated as the durable canonical artifact;
- a squash-path branch is one commit when possible, otherwise its last commit carries the current canonical message;
- every push to a published PR refreshes and verifies the forge metadata;
- evidence supports every concrete claim;
- empty meta, hype, and routine process inventory are removed;
- the description explains behavior and uses references to the diff only when they help orient review;
- no section exceeds the reviewer reading budget without a good reason;
- a reviewer without hidden context can understand purpose, design, verification, and remaining risk; and
- for a substantial change, an independent reviewer can understand the goal and check the safety claim using the commit text and checkout alone.
