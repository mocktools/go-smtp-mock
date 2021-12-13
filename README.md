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
  - [Configuring](#configuring)
  - [Example of usage](#example-of-usage)
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
- Ability to do graceful shutdown of SMTP mock server
- No authentication support
- Zero runtime dependencies
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

- [Configuring](#configuring)
- [Starting-stopping server](#starting-stopping-server)

### Configuring

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
  MsqSizeLimit:                  5,
  

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

  // Custom MAIL FROM blacklisted domain message. Based on defaultQuitMsg by default
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

### Example of usage

You have to create your SMTP mock server using `smtpmock.New()` and `smtpmock.ConfigurationAttr{}` to start interaction with it.

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
  hostAddress, portNumber =: server.PortNumber, "127.0.0.1"

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

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/mocktools/go-smtp-mock>. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct. Please check the [open tikets](https://github.com/mocktools/go-smtp-mock/issues). Be shure to follow Contributor Code of Conduct below and our [Contributing Guidelines](CONTRIBUTING.md).

## License

This golang package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the `smtpmock` project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](CODE_OF_CONDUCT.md).

## Credits

- [The Contributors](https://github.com/mocktools/go-smtp-mock/graphs/contributors) for code and awesome suggestions
- [The Stargazers](https://github.com/mocktools/go-smtp-mock/stargazers) for showing their support

## Versioning

`smtpmock` uses [Semantic Versioning 2.0.0](https://semver.org)
