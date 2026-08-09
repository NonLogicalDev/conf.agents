# Dockerfile Pedantics

## Pins And Generated Blocks

Keep version and digest pins together. When automation owns them, mark the
exact block the updater may replace:

```dockerfile
# BEGIN GENERATED TOOL PIN (managed by manage/update_tool.py; do not edit)
ARG TOOL_VERSION=1.2.3
ARG TOOL_SHA256=...
# END GENERATED TOOL PIN (managed by manage/update_tool.py)
```

The updater should require exactly one ordered marker pair and replace only
that block. Document the updater and its prerequisites beside the image code.

## Build Arguments And Layers

- Name repeated URLs, versions, digests, architectures, and installation paths
  with `ARG` or `ENV` according to whether they remain in the image.
- Keep integrity verification adjacent to the download it verifies.
- Remove temporary installers in the same `RUN` layer that creates them.
- Avoid adding a new stage, init system, or package manager when the existing
  image pattern satisfies the requirement.

## Comments And Validation

Comments should identify ownership, a compatibility constraint worth explaining, or
why a command cannot use the common form. Do not narrate standard Docker
instructions.

Validate every supported architecture. When full image builds are expensive,
run parser or formatter checks first, then the repository's required
image targets for each supported architecture.

For a generated pin block, record the clean file checksum, replace the managed
values with synthetic valid values, run the real updater through its default
path, verify that it restores the canonical values and original checksum, and
run it again to prove that an already current file stays unchanged.
