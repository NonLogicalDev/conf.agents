# Validation Scenarios

## 04 Validate A Managed Pin Refresh With The CLI

### Prompt

Use `$Code_Pedantics`.

The repository command `refresh-toolchain-pin --config <path>` reads the expected `version` and `sha256` from `toolchain.lock`. It rewrites those two lines between `# BEGIN TOOLCHAIN PIN` and `# END TOOLCHAIN PIN` in the selected `toolchain.conf`. The repository includes `testdata/toolchain.conf`, a non-production file with the same marker contract. Unit tests pass. The user asks whether the command is ready. Choose the validation plan. Do not modify the production configuration.

### Expectations

- Copy `testdata/toolchain.conf` to a temporary directory.
- Record the expected `version` and `sha256` from `toolchain.lock`.
- Record all bytes before `# BEGIN TOOLCHAIN PIN` and after `# END TOOLCHAIN PIN`.
- Replace only the temporary file's managed values with a known older `version` and its matching `sha256`.
- Run `refresh-toolchain-pin --config <temporary-path>`.
- Verify the command restored the expected `version` and `sha256`.
- Verify every byte outside the managed block stayed unchanged.
- Run the command again.
- Verify the entire file stayed byte-for-byte unchanged on the second run.
- Do not call the command ready based on unit tests alone.

### Adjacent Valid Case

`refresh-toolchain-pin` can target only the production configuration. It has no dry-run, output-path, temporary-file, or staging mode.

- State that no safe CLI run is available.
- Do not modify production or claim that the actual CLI was validated.
