# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2021-12-07

### Added

- Added ability to assign random SMTP port number by OS as default settings
- Added `Server.PortNumber` field

### Changed

- Updated`Server#Start` method, tests
- Refactored `ConfigurationAttr#assignDefaultValues` method
- Updated `ConfigurationAttr#assignServerDefaultValues` method, tests
- Updated package docs
- Updated linters config

### Removed

- Removed `defaultPortNuber`

## [0.1.2] - 2021-12-03

### Changed

- Updated functions/structures/consts scopes
- Updated linters config
- Updated CircleCI config

### Fixed

- Linters issues

## [0.1.1] - 2021-11-30

### Fixed

- Fixed typos, linter warnings

## [0.1.0] - 2021-11-30

### Added

- First release of `smtpmock`. Thanks [@le0pard](https://github.com/le0pard) for support 🚀
