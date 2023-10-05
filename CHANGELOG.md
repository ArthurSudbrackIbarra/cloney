
# Change Log

All notable changes to this project are documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## Cloney 0.1.1 (Latest) - 2023-10-05
  
### Added

- Added the `docs` command.
 
### Changed
  
- Changed the `dry-run` and `validate` commands to accept a path to a local template repository as the first argument. Before this change, you had to use the `-p, --path` flag to have the same effect.

    ```bash
    # Before
    $ cloney dry-run -p /path/to/template-repo

    # After
    $ cloney dry-run /path/to/template-repo
    ```
 
### Fixed
 
- Fixed a security issue that allowed users to create files and directories outside the scope of the template repository.
- Addressed an issue where the `CLONEY_GIT_TOKEN` environment variable was not being utilized when interacting with private Git repositories.
 
## Cloney 0.1.0 - 2023-10-01
 
This is the first release of Cloney.
 
### Added

- CLI commands: `clone`, `dry-run`, `info`, `start`, `validate`, `version`.
 
### Changed

- No changes.
 
### Fixed

- No fixes.
