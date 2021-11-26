# ![Golang SMTP mock. Mimic SMTP server behaviour for your test environment and even more](https://repository-images.githubusercontent.com/401721985/848bc1dd-fc35-4d78-8bd9-0ac3430270d8)

[![CircleCI](https://circleci.com/gh/mocktools/golang-smtp-mock/tree/master.svg?style=svg)](https://circleci.com/gh/mocktools/golang-smtp-mock/tree/master)
[![GitHub](https://img.shields.io/github/license/mocktools/golang-smtp-mock)](LICENSE.txt)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

Golang SMTP mock. Mimic SMTP server behaviour for your test environment and even more. It's lightweight configurable multithreaded SMTP server package written in Go. It meets the minimum requirements specified by [RFC 2821](https://datatracker.ietf.org/doc/html/rfc2821) & [RFC 5321](https://datatracker.ietf.org/doc/html/rfc5321).

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Code of Conduct](#code-of-conduct)
- [Credits](#credits)
- [Versioning](#versioning)
- [Changelog](CHANGELOG.md)

## Features

- Configurable multithreaded RFC compatible SMTP server
- Implements the minimum command set, responds to commands and adds a valid received header to messages as specified in [RFC 2821](https://datatracker.ietf.org/doc/html/rfc2821) & [RFC 5321](https://datatracker.ietf.org/doc/html/rfc5321)
- Ability to configure behaviour for each SMTP command
- Comes with default settings out of the box, configure only what you need
- Ability to override previous SMTP commands
- Fail fast scenario (closing client session for case when command was inconsistent or failed)
- Mock-server activity logger
- Ability to do graceful shutdown of SMTP mock server
- No authentication support
- Simple and intuitive DSL

## Requirements

Golang 1.15+

## Installation

Install `smtpmock`:

```bash
go get github.com/mocktools/go-smtp-mock
go install -i github.com/mocktools/go-smtp-mock
```

Import `smtpmock` dependency into your code:

```go
package main

import "github.com/mocktools/go-smtp-mock"
```

## Usage

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/mocktools/go-smtp-mock. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct. Please check the [open tikets](https://github.com/mocktools/go-smtp-mock/issues). Be shure to follow Contributor Code of Conduct below and our [Contributing Guidelines](CONTRIBUTING.md).

## License

This golang package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the SmtpMock project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](CODE_OF_CONDUCT.md).

## Credits

- [The Contributors](https://github.com/mocktools/go-smtp-mock/graphs/contributors) for code and awesome suggestions
- [The Stargazers](https://github.com/mocktools/go-smtp-mock/stargazers) for showing their support

## Versioning

SmtpMock uses [Semantic Versioning 2.0.0](https://semver.org)
