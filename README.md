# ![SMTP mock server written on Golang. Mimic any SMTP server behaviour for your test environment with fake SMTP server](https://repository-images.githubusercontent.com/401721985/848bc1dd-fc35-4d78-8bd9-0ac3430270d8)

[![Go Report Card](https://goreportcard.com/badge/github.com/mocktools/go-smtp-mock)](https://goreportcard.com/report/github.com/mocktools/go-smtp-mock)
[![Codecov](https://codecov.io/gh/mocktools/go-smtp-mock/branch/master/graph/badge.svg)](https://codecov.io/gh/mocktools/go-smtp-mock)
[![CircleCI](https://circleci.com/gh/mocktools/go-smtp-mock/tree/master.svg?style=svg)](https://circleci.com/gh/mocktools/go-smtp-mock/tree/master)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/mocktools/go-smtp-mock)](https://github.com/mocktools/go-smtp-mock/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/mocktools/go-smtp-mock)](https://pkg.go.dev/github.com/mocktools/go-smtp-mock)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GitHub](https://img.shields.io/github/license/mocktools/go-smtp-mock)](LICENSE.txt)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

`smtpmock` is lightweight configurable multithreaded fake SMTP server written in Go. It meets the minimum requirements specified by [RFC 2821](https://datatracker.ietf.org/doc/html/rfc2821) & [RFC 5321](https://datatracker.ietf.org/doc/html/rfc5321). Allows to mimic any SMTP server behaviour for your test environment and even more 🚀

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Inside of Golang ecosystem](#inside-of-golang-ecosystem)
    - [Configuring](#configuring)
    - [Manipulation with server](#manipulation-with-server)
  - [Inside of Ruby ecosystem](#inside-of-ruby-ecosystem)
    - [Configuring](#configuring)
    - [Manipulation with server](#manipulation-with-server)
  - [Inside of any ecosystem](#inside-of-any-ecosystem)
    - [Configuring with command line arguments](#configuring-with-command-line-arguments)
    - [Other options](#other-options)
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
- Fail fast scenario (ability to close client session for case when command was inconsistent or failed)
- Mock-server activity logger
- Ability to do graceful/force shutdown of SMTP mock server
- No authentication support
- Zero runtime dependencies
- Simple and intuitive DSL
- Ability to run server as binary with command line arguments

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

- [Inside of Golang ecosystem](#inside-of-golang-ecosystem)
  - [Configuring](#configuring)
  - [Manipulation with server](#manipulation-with-server)
- [Inside of Ruby ecosystem](#inside-of-ruby-ecosystem)
  - [Configuring](#configuring)
  - [Manipulation with server](#manipulation-with-server)
- [Inside of any ecosystem](#inside-of-any-ecosystem)
  - [Configuring with command line arguments](#configuring-with-command-line-arguments)
  - [Other options](#other-options)

### Inside of Golang ecosystem

You have to create your SMTP mock server using `smtpmock.New()` and `smtpmock.ConfigurationAttr{}` to start interaction with it.

#### Configuring

`smtpmock` is SMTP server for test environment with configurable behaviour. It comes with default settings out of the box. But you can override any default behaviour if you need.

```go
smtpmock.ConfigurationAttr{

  // Customizing server behaviour
  // ---------------------------------------------------------------------
  // Host address where smtpmock will run, it's equal to "127.0.0.1" by default
  HostAddress:                   "[::]",

  // Port number on which the server will bind. If it not specified, it will be
  // assigned dynamically after server.Start() by default
  PortNumber:                    2525,

  // Enables/disables log to stdout. It's equal to false by default
  LogToStdout:                   true,

  // Enables/disables log server activity. It's equal to false by default
  LogServerActivity:             true,

  // Ability to specify session timeout. It's equal to 30 seconds by default
  SessionTimeout:                42,

  // Ability to specify graceful shutdown timeout. It's equal to 1 second by default
  ShutdownTimeout:               5,


  // Customizing SMTP command handlers behaviour
  // ---------------------------------------------------------------------
  // Ability to configure fail fast scenario. It means that server will
  // close client session for case when command was inconsistent or failed.
  // It's equal to false by default
  IsCmdFailFast:                 true,

  // Ability to specify blacklisted HELO domains. It's equal to empty []string
  BlacklistedHeloDomains:        []string{"example1.com", "example2.com", "localhost"},

  // Ability to specify blacklisted MAIL FROM emails. It's equal to empty []string
  BlacklistedMailfromEmails:     []string{"bot@olo.com", "robot@molo.com"},

  // Ability to specify blacklisted RCPT TO emails. It's equal to empty []string
  BlacklistedRcpttoEmails:       []string{"blacklisted@olo.com", "blacklisted@molo.com"},

  // Ability to specify not registered (non-existent) RCPT TO emails.
  // It's equal to empty []string
  NotRegisteredEmails:           []string{"nobody@olo.com", "non-existent@email.com"},

  // Ability to specify message body size limit. It's equal to 10485760 bytes (10MB) by default
  MsgSizeLimit:                  5,
  

  // Customazing SMTP command handler messages context
  // ---------------------------------------------------------------------
  // Custom server greeting message. Base on defaultGreetingMsg by default
  MsgGreeting:                   "msgGreeting",

  // Custom invalid command message. Based on defaultInvalidCmdMsg by default
  MsgInvalidCmd:                 "msgInvalidCmd",

  // Custom invalid command HELO sequence message.
  // Based on defaultInvalidCmdHeloSequenceMsg by default
  MsgInvalidCmdHeloSequence:     "msgInvalidCmdHeloSequence",

  // Custom invalid command HELO argument message.
  // Based on defaultInvalidCmdHeloArgMsg by default
  MsgInvalidCmdHeloArg:          "msgInvalidCmdHeloArg",

  // Custom HELO blacklisted domain message. Based on defaultQuitMsg by default
  MsgHeloBlacklistedDomain:      "msgHeloBlacklistedDomain",

  // Custom HELO received message. Based on defaultReceivedMsg by default
  MsgHeloReceived:               "msgHeloReceived",

  // Custom invalid command MAIL FROM sequence message.
  // Based on defaultInvalidCmdMailfromSequenceMsg by default
  MsgInvalidCmdMailfromSequence: "msgInvalidCmdMailfromSequence",

  // Custom invalid command MAIL FROM argument message.
  // Based on defaultInvalidCmdMailfromArgMsg by default
  MsgInvalidCmdMailfromArg:      "msgInvalidCmdMailfromArg",

  // Custom MAIL FROM blacklisted email message. Based on defaultQuitMsg by default
  MsgMailfromBlacklistedEmail:   "msgMailfromBlacklistedEmail",

  // Custom MAIL FROM received message. Based on defaultReceivedMsg by default
  MsgMailfromReceived:           "msgMailfromReceived",

  // Custom invalid command RCPT TO sequence message.
  // Based on defaultInvalidCmdRcpttoSequenceMsg by default
  MsgInvalidCmdRcpttoSequence:   "msgInvalidCmdRcpttoSequence",

  // Custom invalid command RCPT TO argument message.
  // Based on defaultInvalidCmdRcpttoArgMsg by default
  MsgInvalidCmdRcpttoArg:        "msgInvalidCmdRcpttoArg",

  // Custom RCPT TO not registered email message.
  // Based on defaultNotRegistredRcpttoEmailMsg by default
  MsgRcpttoNotRegisteredEmail:   "msgRcpttoNotRegisteredEmail",

  // Custom RCPT TO blacklisted email message. Based on defaultQuitMsg by default
  MsgRcpttoBlacklistedEmail:     "msgRcpttoBlacklistedEmail",

  // Custom RCPT TO received message. Based on defaultReceivedMsg by default
  MsgRcpttoReceived:             "msgRcpttoReceived",

  // Custom invalid command DATA sequence message.
  // Based on defaultInvalidCmdDataSequenceMsg by default
  MsgInvalidCmdDataSequence:     "msgInvalidCmdDataSequence",

  // Custom DATA received message. Based on defaultReadyForReceiveMsg by default
  MsgDataReceived:               "msgDataReceived",

  // Custom size is too big message. Based on defaultMsgSizeIsTooBigMsg by default
  MsgMsgSizeIsTooBig:            "msgMsgSizeIsTooBig",

  // Custom received message body message. Based on defaultReceivedMsg by default
  MsgMsgReceived:                "msgMsgReceived",

  // Custom quit command message. Based on defaultQuitMsg by default
  MsgQuitCmd:                    "msgQuitCmd",
}
```

#### Manipulation with server

```go
package main

import (
  "fmt"
  "net"
  "net/smtp"

  "github.com/mocktools/go-smtp-mock"
)

func main() {
  // You can pass empty smtpmock.ConfigurationAttr{}. It means that smtpmock will use default settings
  server := smtpmock.New(smtpmock.ConfigurationAttr{
    LogToStdout:       true,
    LogServerActivity: true,
  })

  // To start server use Start() method
  if err := server.Start(); err != nil {
    fmt.Println(err)
  }

  // Server's port will be assigned dynamically after server.Start()
  // for case when portNumber wasn't specified
  hostAddress, portNumber := "127.0.0.1", server.PortNumber

  // Possible SMTP-client stuff for iteration with mock server
  address := fmt.Sprintf("%s:%d", hostAddress, portNumber)
  timeout := time.Duration(2) * time.Second

  connection, _ := net.DialTimeout("tcp", address, timeout)
  client, _ := smtp.NewClient(connection, hostAddress)
  client.Hello("example.com")
  client.Quit()
  client.Close()

  // To stop the server use Stop() method. Please note, smtpmock uses graceful shutdown.
  // It means that smtpmock will end all sessions after client responses or by session
  // timeouts immediately.
  if err := server.Stop(); err != nil {
    fmt.Println(err)
  }
}
```

Code from example above will produce next output to stdout:

```
INFO: 2021/11/30 22:07:30.554827 SMTP mock server started on port: 2525
INFO: 2021/11/30 22:07:30.554961 SMTP session started
INFO: 2021/11/30 22:07:30.554998 SMTP response: 220 Welcome
INFO: 2021/11/30 22:07:30.555059 SMTP request: EHLO example.com
INFO: 2021/11/30 22:07:30.555648 SMTP response: 250 Received
INFO: 2021/11/30 22:07:30.555686 SMTP request: QUIT
INFO: 2021/11/30 22:07:30.555722 SMTP response: 221 Closing connection
INFO: 2021/11/30 22:07:30.555732 SMTP session finished
WARNING: 2021/11/30 22:07:30.555801 SMTP mock server is in the shutdown mode and won't accept new connections
INFO: 2021/11/30 22:07:30.555808 SMTP mock server was stopped successfully
```

### Inside of Ruby ecosystem

In Ruby ecosystem `smtpmock` is available as [`smtp_mock`](https://github.com/mocktools/ruby-smtp-mock) gem. It's flexible Ruby wrapper over `smtpmock` binary.

#### Configuring

It comes with default settings out of the box. List of all [available server options](https://github.com/mocktools/ruby-smtp-mock#available-server-options). Configure only what you need, for example:

```ruby
configuration = { not_registered_emails: %w[user@olo.com user@molo.com] }
```

#### Manipulation with server

First, you should install `smtp_mock` gem and `smtpmock` as system dependency:

```bash
gem install smtp_mock
bundle exec smtp_mock -i ~
```

Now, you can create and interact with your `smtpmock` instance natively from Ruby ecosystem:

```ruby
require 'smtp_mock'

smtp_mock_server = SmtpMock.start_server(**configuration)

# returns current smtp mock server port
smtp_mock_server.port # => 55640

# interface for force shutdown current smtp mock server
smtp_mock_server.stop! # => true
```

### Inside of any ecosystem

You can use `smtpmock` as binary. Just download the pre-compiled binary from the [releases page](https://github.com/mocktools/go-smtp-mock/releases) and copy them to the desired location. For start server run command with needed arguments. You can use our bash script for automation this process like in the example below:

```bash
curl -sL https://raw.githubusercontent.com/mocktools/go-smtp-mock/master/script/download.sh | bash
./smtpmock -port=2525 -log
```

#### Configuring with command line arguments

`smtpmock` configuration is available as command line arguments specified in the list below:

| Flag description | Example of usage |
| --- | --- |
| `-host` - host address where smtpmock will run. It's equal to `127.0.0.1` by default | `-host=localhost` |
| `-port` - server port number. If not specified it will be assigned dynamically | `-port=2525` |
| `-log` - enables log server activity. Disabled by default | `-log` |
| `-sessionTimeout` - session timeout in seconds. It's equal to 30 seconds by default | `-sessionTimeout=60` |
| `-shutdownTimeout` - graceful shutdown timeout in seconds. It's equal to 1 second by default | `-shutdownTimeout=5` |
| `-failFast` - enables fail fast scenario. Disabled by default | `-failFast` |
| `-blacklistedHeloDomains` - blacklisted `HELO` domains, separated by commas | `-blacklistedHeloDomains="example1.com,example2.com"` |
| `-blacklistedMailfromEmails` - blacklisted `MAIL FROM` emails, separated by commas | `-blacklistedMailfromEmails="a@example1.com,b@example2.com"` |
| `-blacklistedRcpttoEmails` - blacklisted `RCPT TO` emails, separated by commas | `-blacklistedRcpttoEmails="a@example1.com,b@example2.com"` |
| `-notRegisteredEmails` - not registered (non-existent) `RCPT TO` emails, separated by commas | `-notRegisteredEmails="a@example1.com,b@example2.com"` |
| `-msgSizeLimit` - message body size limit in bytes. It's equal to `10485760` bytes | `-msgSizeLimit=42` |
| `-msgGreeting` - custom server greeting message | `-msgGreeting="Greeting message"` |
| `-msgInvalidCmd` - custom invalid command message | `-msgInvalidCmd="Invalid command message"` |
| `-msgInvalidCmdHeloSequence` - custom invalid command `HELO` sequence message | `-msgInvalidCmdHeloSequence="Invalid command HELO sequence message"` |
| `-msgInvalidCmdHeloArg` - custom invalid command `HELO` argument message | `-msgInvalidCmdHeloArg="Invalid command HELO argument message"` |
| `-msgHeloBlacklistedDomain` - custom `HELO` blacklisted domain message | `-msgHeloBlacklistedDomain="Blacklisted domain message"` |
| `-msgHeloReceived` - custom `HELO` received message | `-msgHeloReceived="HELO received message"` |
| `-msgInvalidCmdMailfromSequence` - custom invalid command `MAIL FROM` sequence message | `-msgInvalidCmdMailfromSequence="Invalid command MAIL FROM sequence message"` |
| `-msgInvalidCmdMailfromArg` - custom invalid command `MAIL FROM` argument message | `-msgInvalidCmdMailfromArg="Invalid command MAIL FROM argument message"` |
| `-msgMailfromBlacklistedEmail` - custom `MAIL FROM` blacklisted email message | `-msgMailfromBlacklistedEmail="Blacklisted email message"` |
| `-msgMailfromReceived`- custom `MAIL FROM` received message | `-msgMailfromReceived="MAIL FROM received message"` |
| `-msgInvalidCmdRcpttoSequence` - custom invalid command `RCPT TO` sequence message | `-msgInvalidCmdRcpttoSequence="Invalid command RCPT TO sequence message"` |
| `-msgInvalidCmdRcpttoArg` - custom invalid command `RCPT TO` argument message | `-msgInvalidCmdRcpttoArg="Invalid command RCPT TO argument message"` |
| `-msgRcpttoNotRegisteredEmail` - custom `RCPT TO` not registered email message | `-msgRcpttoNotRegisteredEmail="Not registered email message"` |
| `-msgRcpttoBlacklistedEmail` - custom `RCPT TO` blacklisted email message | `-msgRcpttoBlacklistedEmail="Blacklisted email message"` |
| `-msgRcpttoReceived` - custom `RCPT TO` received message | `-msgRcpttoReceived="RCPT TO received message"` |
| `-msgInvalidCmdDataSequence` - custom invalid command `DATA` sequence message | `-msgInvalidCmdDataSequence="Invalid command DATA sequence message"` |
| `-msgDataReceived` - custom `DATA` received message | `-msgDataReceived="DATA received message"` |
| `-msgMsgSizeIsTooBig` - custom size is too big message | `-msgMsgSizeIsTooBig="Message size is too big"` |
| `-msgMsgReceived` - custom received message body message | `-msgMsgReceived="Message has been received"` |
| `-msgQuitCmd` - custom quit command message | `-msgQuitCmd="Quit command message"` |

#### Other options

Available not configuration `smtpmock` options:

| Flag description | Example of usage |
| --- | --- |
| `-v` - Just prints current `smtpmock` binary build data (version, commit, datetime). Doesn't run the server. | `-v` |

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/mocktools/go-smtp-mock>. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct. Please check the [open tickets](https://github.com/mocktools/go-smtp-mock/issues). Be sure to follow Contributor Code of Conduct below and our [Contributing Guidelines](CONTRIBUTING.md).

## License

This golang package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the `smtpmock` project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](CODE_OF_CONDUCT.md).

## Credits

- [The Contributors](https://github.com/mocktools/go-smtp-mock/graphs/contributors) for code and awesome suggestions
- [The Stargazers](https://github.com/mocktools/go-smtp-mock/stargazers) for showing their support

## Versioning

`smtpmock` uses [Semantic Versioning 2.0.0](https://semver.org)
