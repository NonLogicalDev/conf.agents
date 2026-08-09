# Dockerfile Scenarios

## 23 Preserve Docker Build Ownership

### Prompt

Use `$Code_Pedantics`.

A Dockerfile keeps a tool version and digest in an unmarked block that an updater rewrites with a broad regular expression. The version is duplicated in an `ENV` even though it is needed only during build, while a runtime path is kept as an `ARG`. The download layer verifies a digest but leaves its installer behind. The change was tested only with an amd64 build even though the image declares amd64 and arm64 support. Choose the review result. Do not modify files.

### Expectations

- Keep the version and digest together inside one exact generated marker pair, and make the updater replace only that owned block.
- Use `ARG` for values needed only during the build and `ENV` only for values that remain in the final image.
- Keep integrity verification beside the download and remove temporary installers in the same `RUN` layer that creates them.
- Validate every declared architecture, not just the developer's host architecture.
- For the generated block, require a clean checksum, synthetic valid values, a real updater run through its default path, restored canonical output, and a second no-op run.

### Adjacent Valid Case

The image intentionally supports only one architecture, and the updater owns a single ordered marker pair with a checked no-op path.

- Validate the one declared architecture instead of inventing unsupported targets.
- Keep the generated block explicit and verify that its updater works.
