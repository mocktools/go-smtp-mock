package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"syscall"
	"testing"

	version "github.com/mocktools/go-smtp-mock/cmd/version"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Run("when error not happened", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		signals <- syscall.SIGINT
		main()
	})

	t.Run("when error happened", func(t *testing.T) {
		defer func() { logFatalf = log.Fatalf }()
		os.Args = []string{os.Args[0], "-host=a"}
		logMock := new(logMock)
		logFatalf = logMock.Fatalf
		errorInterface := []interface{}{errors.New("Failed to start SMTP mock server on port: 0")}
		logMock.On("Fatalf", "%s\n", errorInterface).Once().Return(nil)
		main()
	})
}

func TestRun(t *testing.T) {
	path := "some-path-to-the-program"

	t.Run("when command line argument error", func(t *testing.T) {
		assert.Error(t, run([]string{path, "-port=a"}, flag.ContinueOnError))
	})

	t.Run("when server starting error", func(t *testing.T) {
		assert.Error(t, run([]string{path, "-host=a"}))
	})

	t.Run("when server was started successfully, interrupt signal (exit 2) received", func(t *testing.T) {
		signals <- syscall.SIGINT

		assert.NoError(t, run([]string{path}))
	})

	t.Run("when server was started successfully, quit signal (exit 3) received", func(t *testing.T) {
		signals <- syscall.SIGQUIT

		assert.NoError(t, run([]string{path}))
	})

	t.Run("when server was started successfully, terminated signal (exit 15) received", func(t *testing.T) {
		signals <- syscall.SIGTERM

		assert.NoError(t, run([]string{path}))
	})

	t.Run("when version flag passed", func(t *testing.T) {
		assert.NoError(t, run([]string{path, "-v"}))
	})
}

func TestToSlice(t *testing.T) {
	t.Run("converts string separated by commas to slice of strings", func(t *testing.T) {
		assert.Equal(t, []string{"a", "b"}, toSlice("a,b"))
	})
}

func TestPrintVersionData(t *testing.T) {
	t.Run("", func(t *testing.T) {
		bytesBuffer := new(bytes.Buffer)
		printVersionData(bytesBuffer)
		ver := "smtpmock: " + version.Version + "\n"
		commit := "commit: " + version.GitCommit + "\n"
		builtAt := "built at: " + version.BuildTime + "\n"
		versionData := ver + commit + builtAt
		assert.Equal(t, versionData, bytesBuffer.String())
	})
}

func TestAttrFromCommandLine(t *testing.T) {
	t.Run("when known flags found creates pointer to ConfigurationAttr based on passed command line arguments", func(t *testing.T) {
		hostAddress := "0"
		portNumber := 42
		sessionTimeout := 12
		shutdownTimeout := 5
		blacklistedHeloDomains := "a.com,b.com"
		blacklistedMailfromEmails := "a@a.com,b@b.com"
		blacklistedRcpttoEmails := "c@a.com,d@b.com"
		notRegisteredEmails := "non-existent@a.com"
		responseDelayHelo := 1
		responseDelayMailfrom := 2
		responseDelayRcptto := 3
		responseDelayData := 4
		responseDelayMessage := 5
		responseDelayRset := 6
		responseDelayQuit := 7
		msgSizeLimit := 1000
		msgGreeting := "msgGreeting"
		msgInvalidCmd := "msgInvalidCmd"
		msgInvalidCmdHeloSequence := "msgInvalidCmdHeloSequence"
		msgInvalidCmdHeloArg := "msgInvalidCmdHeloArg"
		msgHeloBlacklistedDomain := "msgHeloBlacklistedDomain"
		msgHeloReceived := "msgHeloReceived"
		msgInvalidCmdMailfromSequence := "msgInvalidCmdMailfromSequence"
		msgInvalidCmdMailfromArg := "msgInvalidCmdMailfromArg"
		msgMailfromBlacklistedEmail := "msgMailfromBlacklistedEmail"
		msgMailfromReceived := "msgMailfromReceived"
		msgInvalidCmdRcpttoSequence := "msgInvalidCmdRcpttoSequence"
		msgInvalidCmdRcpttoArg := "msgInvalidCmdRcpttoArg"
		msgRcpttoNotRegisteredEmail := "msgRcpttoNotRegisteredEmail"
		msgRcpttoBlacklistedEmail := "msgRcpttoBlacklistedEmail"
		msgRcpttoReceived := "msgRcpttoReceived"
		msgInvalidCmdDataSequence := "msgInvalidCmdDataSequence"
		msgDataReceived := "msgDataReceived"
		msgMsgSizeIsTooBig := "msgMsgSizeIsTooBig"
		msgMsgReceived := "msgMsgReceived"
		msgRsetReceived := "msgRsetReceived"
		msgQuitCmd := "msgQuitCmd"
		ver, configAttr, err := attrFromCommandLine(
			[]string{
				"some-path-to-the-program",
				"-v",
				"-host=" + hostAddress,
				"-port=" + strconv.Itoa(portNumber),
				"-log",
				"-sessionTimeout=" + strconv.Itoa(sessionTimeout),
				"-shutdownTimeout=" + strconv.Itoa(shutdownTimeout),
				"-failFast",
				"-multipleMessageReceiving",
				"-blacklistedHeloDomains=" + blacklistedHeloDomains,
				"-blacklistedMailfromEmails=" + blacklistedMailfromEmails,
				"-blacklistedRcpttoEmails=" + blacklistedRcpttoEmails,
				"-notRegisteredEmails=" + notRegisteredEmails,
				"-responseDelayHelo=" + strconv.Itoa(responseDelayHelo),
				"-responseDelayMailfrom=" + strconv.Itoa(responseDelayMailfrom),
				"-responseDelayRcptto=" + strconv.Itoa(responseDelayRcptto),
				"-responseDelayData=" + strconv.Itoa(responseDelayData),
				"-responseDelayMessage=" + strconv.Itoa(responseDelayMessage),
				"-responseDelayRset=" + strconv.Itoa(responseDelayRset),
				"-responseDelayQuit=" + strconv.Itoa(responseDelayQuit),
				"-msgSizeLimit=" + strconv.Itoa(msgSizeLimit),
				"-msgGreeting=" + msgGreeting,
				"-msgInvalidCmd=" + msgInvalidCmd,
				"-msgInvalidCmdHeloSequence=" + msgInvalidCmdHeloSequence,
				"-msgInvalidCmdHeloArg=" + msgInvalidCmdHeloArg,
				"-msgHeloBlacklistedDomain=" + msgHeloBlacklistedDomain,
				"-msgHeloReceived=" + msgHeloReceived,
				"-msgInvalidCmdMailfromSequence=" + msgInvalidCmdMailfromSequence,
				"-msgInvalidCmdMailfromArg=" + msgInvalidCmdMailfromArg,
				"-msgMailfromBlacklistedEmail=" + msgMailfromBlacklistedEmail,
				"-msgMailfromReceived=" + msgMailfromReceived,
				"-msgInvalidCmdRcpttoSequence=" + msgInvalidCmdRcpttoSequence,
				"-msgInvalidCmdRcpttoArg=" + msgInvalidCmdRcpttoArg,
				"-msgRcpttoNotRegisteredEmail=" + msgRcpttoNotRegisteredEmail,
				"-msgRcpttoBlacklistedEmail=" + msgRcpttoBlacklistedEmail,
				"-msgRcpttoReceived=" + msgRcpttoReceived,
				"-msgInvalidCmdDataSequence=" + msgInvalidCmdDataSequence,
				"-msgDataReceived=" + msgDataReceived,
				"-msgMsgSizeIsTooBig=" + msgMsgSizeIsTooBig,
				"-msgMsgReceived=" + msgMsgReceived,
				"-msgRsetReceived=" + msgRsetReceived,
				"-msgQuitCmd=" + msgQuitCmd,
			},
		)

		assert.True(t, ver)
		assert.Equal(t, hostAddress, configAttr.HostAddress)
		assert.Equal(t, portNumber, configAttr.PortNumber)
		assert.True(t, configAttr.LogToStdout)
		assert.True(t, configAttr.LogServerActivity)
		assert.Equal(t, sessionTimeout, configAttr.SessionTimeout)
		assert.Equal(t, shutdownTimeout, configAttr.ShutdownTimeout)
		assert.True(t, configAttr.IsCmdFailFast)
		assert.True(t, configAttr.MultipleMessageReceiving)
		assert.Equal(t, toSlice(blacklistedHeloDomains), configAttr.BlacklistedHeloDomains)
		assert.Equal(t, toSlice(blacklistedMailfromEmails), configAttr.BlacklistedMailfromEmails)
		assert.Equal(t, toSlice(blacklistedRcpttoEmails), configAttr.BlacklistedRcpttoEmails)
		assert.Equal(t, toSlice(notRegisteredEmails), configAttr.NotRegisteredEmails)
		assert.Equal(t, responseDelayHelo, configAttr.ResponseDelayHelo)
		assert.Equal(t, responseDelayMailfrom, configAttr.ResponseDelayMailfrom)
		assert.Equal(t, responseDelayRcptto, configAttr.ResponseDelayRcptto)
		assert.Equal(t, responseDelayData, configAttr.ResponseDelayData)
		assert.Equal(t, responseDelayMessage, configAttr.ResponseDelayMessage)
		assert.Equal(t, responseDelayRset, configAttr.ResponseDelayRset)
		assert.Equal(t, responseDelayQuit, configAttr.ResponseDelayQuit)
		assert.Equal(t, msgSizeLimit, configAttr.MsgSizeLimit)
		assert.Equal(t, msgGreeting, configAttr.MsgGreeting)
		assert.Equal(t, msgInvalidCmd, configAttr.MsgInvalidCmd)
		assert.Equal(t, msgInvalidCmdHeloSequence, configAttr.MsgInvalidCmdHeloSequence)
		assert.Equal(t, msgInvalidCmdHeloArg, configAttr.MsgInvalidCmdHeloArg)
		assert.Equal(t, msgHeloBlacklistedDomain, configAttr.MsgHeloBlacklistedDomain)
		assert.Equal(t, msgHeloReceived, configAttr.MsgHeloReceived)
		assert.Equal(t, msgInvalidCmdMailfromSequence, configAttr.MsgInvalidCmdMailfromSequence)
		assert.Equal(t, msgInvalidCmdMailfromArg, configAttr.MsgInvalidCmdMailfromArg)
		assert.Equal(t, msgMailfromBlacklistedEmail, configAttr.MsgMailfromBlacklistedEmail)
		assert.Equal(t, msgMailfromReceived, configAttr.MsgMailfromReceived)
		assert.Equal(t, msgInvalidCmdRcpttoSequence, configAttr.MsgInvalidCmdRcpttoSequence)
		assert.Equal(t, msgInvalidCmdRcpttoArg, configAttr.MsgInvalidCmdRcpttoArg)
		assert.Equal(t, msgRcpttoNotRegisteredEmail, configAttr.MsgRcpttoNotRegisteredEmail)
		assert.Equal(t, msgRcpttoBlacklistedEmail, configAttr.MsgRcpttoBlacklistedEmail)
		assert.Equal(t, msgRcpttoReceived, configAttr.MsgRcpttoReceived)
		assert.Equal(t, msgInvalidCmdDataSequence, configAttr.MsgInvalidCmdDataSequence)
		assert.Equal(t, msgDataReceived, configAttr.MsgDataReceived)
		assert.Equal(t, msgMsgSizeIsTooBig, configAttr.MsgMsgSizeIsTooBig)
		assert.Equal(t, msgMsgReceived, configAttr.MsgMsgReceived)
		assert.Equal(t, msgRsetReceived, configAttr.MsgRsetReceived)
		assert.Equal(t, msgQuitCmd, configAttr.MsgQuitCmd)
		assert.NoError(t, err)
	})

	t.Run("when unknown flags found sends exit signal", func(t *testing.T) {
		ver, configAttr, err := attrFromCommandLine([]string{"some-path-to-the-program", "-notKnownFlag"}, flag.ContinueOnError)

		assert.False(t, ver)
		assert.Nil(t, configAttr)
		assert.Error(t, err)
	})
}
