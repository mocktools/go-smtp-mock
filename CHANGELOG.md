# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.1] - 2022-01-02

### Added

- Added bash script for downloading latest release

### Updated

- Updated release binary package signature
- Updated package documentation

## [1.4.0] - 2021-12-20

### Added

- Implemented ability to do force stop by timeout

### Updated

- Updated `configuration`, tests
- Updated `server`, tests
- Updated `main`, tests
- Updated consts, package documentation

## [1.3.5] - 2021-12-16

### Updated

- Updated CircleCI config
- Updated goreleaser config

## [1.3.4] - 2021-12-16

### Updated

- Updated CircleCI config

## [1.3.3] - 2021-12-16

### Updated

- Updated CircleCI config

## [1.3.2] - 2021-12-16

### Updated

- Updated CircleCI config

## [1.3.1] - 2021-12-16

### Added

- Added goreleaser config

## [1.3.0] - 2021-12-16

### Added

- Added ability to run smtpmock as service
- Implemented package main, tests

### Fixed

- Fixed documentation issues. Thanks [@vpakhuchyi](https://github.com/vpakhuchyi) for report and PR.
- Fixed `MsgSizeLimit`, `msgSizeLimit` typo in fields naming. Thanks [@vanyavasylyshyn](https://github.com/vanyavasylyshyn) for report.
- Fixed project gihub templates

### Updated

- Updated CircleCI config
- Updated package documentation

## [1.2.0] - 2021-12-13

### Added

- Added ability to use localhost as valid `HELO` domain. Thanks [@lesichkovm](https://github.com/lesichkovm) for report.

### Changed

- Updated `handlerHelo#heloDomain`, tests
- Updated consts
- Updated package docs

## [1.1.0] - 2021-12-11

### Changed

- Updated default negative SMTP command responses follows to RFC
- Updated `ConfigurationAttr` methods, tests

## [1.0.1] - 2021-12-10

### Fixed

- Fixed `ConfigurationAttr` unexported fields issue

### Changed

- Updated package documentation

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
