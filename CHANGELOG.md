# Change Log

This document lists all significant changes to the Cloney project, following [Keep a Changelog](http://keepachangelog.com/) and adhering to [Semantic Versioning](http://semver.org/).

## (Major) Cloney 1.0.0 - 2023-11-XX

### Added

- Introducing the `joinDoubleQuote` and `joinSingleQuote` template functions. They join a list with a given separator and surround each item with double/single quotes.
- Error feedback has been added to inform users when they use an unknown command or unsupported flags.

### Changed

- Enhanced template repository handling:
  - Files and directories starting with `__` (double underscore) are now ignored, allowing them to be used for internal purposes within the template repository without being included in the final output. For example, you can use these files to define templates that are reused in other files, ensuring they are not part of the output.
  - Template definitions can now be shared between multiple files, promoting reusability.
  - Automatic exclusion of known files and directories, such as `node_modules` and `.venv`, to prevent unnecessary processing.

### Fixed

- Resolved Windows-specific issues:
  - Prevented occasional crashes during the `dry-run` and `clone` commands due to path problems.
  - Fixed the issue where backslashes (`\`) were sometimes replaced with forward slashes (`/`) on Windows.
- Resolved an issue where command outputs were incorrectly directed to `stderr` instead of `stdout`.

## (Minor) Cloney 0.2.0 - 2023-10-05

### Added

- Introduced the `docs` command.

### Changed

- Simplified command syntax for the `dry-run` and `validate` commands by accepting a path to a local template repository as the first argument.

  ```bash
  # Before
  $ cloney dry-run -p /path/to/template-repo

  # After
  $ cloney dry-run /path/to/template-repo
  ```

### Fixed

- Addressed security and functionality concerns:
  - Fixed a security issue that allowed users to create files and directories outside the template repository's scope.
  - Corrected the handling of the `CLONEY_GIT_TOKEN` environment variable for interacting with private Git repositories.

## Cloney 0.1.0 - 2023-10-01

This marks the initial release of Cloney.

### Added

- Introduced a set of CLI commands: `clone`, `dry-run`, `info`, `start`, `validate`, and `version`.

### Changed

- No changes made in this release.

### Fixed

- No bug fixes were required for this release.
