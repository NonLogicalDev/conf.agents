# Asking good questions

Read this reference when asking the user to choose, preparing a technical question, asking someone for help, checking your understanding, or deciding whom to interrupt.

The advice on asking for help draws from [Julia Evans: How to ask good questions](https://jvns.ca/blog/good-questions/).

## Ask for one useful decision

Ask only when the answer is needed to continue or would materially improve the result. Make the decision clear before presenting the choices:

- Say what the choices already include when it helps the reader distinguish the options.
- Describe each choice by what the user will see or be able to do. Explain the practical consequence in familiar words.
- Recommend a sensible default and briefly explain why it fits. If no choice has a clear advantage, say what the choice depends on.
- Include only the context needed to choose. Give question titles and options the same care as the surrounding prose.
- Omit process labels, ceremony, and alternatives that do not represent a useful choice. Put answers and status updates in normal chat.

The Flow examples below show how shared context makes two separate decisions easier to answer. Their wording, numbering, and length are not a required template.

### Example: choose an explanation format

**Which default explanation format should Flow use?** Both choices include semantic pseudocode and source references; this decides how much surrounding prose appears.

- **Brief review explanation (Recommended):** A short Why/What, a diff with context, source references, and only consequences that affect the review. This keeps routine reviews quick to scan.
- **More surrounding prose:** The same material with fuller explanations around the pseudocode and references. This gives readers more guidance to work through.

### Example: choose what a source link does

**When a source reference is clicked, should Flow jump to the cited lines in its Files section or expand a code snippet under the explanation?** Both choices stay on the same Delta page.

- **Jump to Files (Recommended):** Highlight the cited old or new lines in the existing diff so the reviewer can see the surrounding code.
- **Expand a snippet:** Show the cited code under the explanation so the reviewer can read it there.

## Make the question easy to answer

A small amount of research can make your question much clearer. Use what you learned to give the other person enough context:

1. Start with the outcome you are trying to achieve.
2. State your current understanding and the relevant context.
3. Explain what you checked and where you are stuck.

Ask a concrete question rather than requesting a broad explanation:

- Say what you believe so the other person can correct it.
- Separate what you observed from what you inferred.
- Name any unfamiliar term or missing concept.
- Ask for a specific fact, clarification, decision, or next step.
- Include enough evidence to make the question understandable on its own.

A question can put the setup in a short list:

- I am trying to achieve [outcome].
- My current understanding is [context].
- I checked [evidence].

Is [specific understanding] correct, or should I look at [specific alternative]?

## Choose whom and when to ask

Choose a person with relevant knowledge or ownership. Before asking, consider:

- Are they available?
- How long will an answer take?
- Is the expected benefit worth the interruption?

## Leave room to learn

Imperfect questions are normal. Asking for clarification is part of learning, and a thoughtful question can make hidden assumptions visible to everyone.

Preparation should make collaboration easier. Do not turn it into a test, discourage questions, or require someone to exhaust every possible alternative before asking for help.

Share useful answers when they will help other people too.
