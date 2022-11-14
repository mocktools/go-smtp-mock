package smtpmock

import (
	"io"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

// Returns log message regex based on log level and message context
func loggerMessageRegex(logLevel, logMessage string) *regexp.Regexp {
	regex, _ := newRegex(logLevel + `: \d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}\.\d{6} ` + logMessage)
	return regex
}

// Creates configuration with default settings
func createConfiguration() *configuration {
	return newConfiguration(ConfigurationAttr{})
}

// Creates not empty message
func createNotEmptyMessage() *Message {
	return &Message{
		heloRequest:           "a",
		heloResponse:          "b",
		mailfromRequest:       "c",
		mailfromResponse:      "d",
		rcpttoRequestResponse: [][]string{[]string{"request", "response"}},
		dataRequest:           "c",
		dataResponse:          "d",
		msgRequest:            "a",
		msgResponse:           "b",
		rsetRequest:           "a",
		rsetResponse:          "b",
		helo:                  true,
		mailfrom:              true,
		rcptto:                true,
		data:                  true,
		msg:                   true,
		rset:                  true,
	}
}

// Creates array of bytes with message body
func messageBody(from, to string) []byte {
	return []byte(
		strings.Join(
			[]string{
				"From: " + from,
				"To: " + to,
				"Subject: Test message for: " + to,
				"Content-Type: text/html; charset=utf-8;",
				"Test message context",
			},
			"\r\n",
		),
	)
}

// Runs full smtp flow
func runFullFlow(client *smtp.Client) error {
	var err error
	var wc io.WriteCloser

	sender, receiver1, receiver2, receiver3 := "user@molo.com", "user1@olo.com", "user2@olo.com", "user3@olo.com"

	if err = client.Mail(sender); err != nil {
		return err
	}
	if err = client.Rcpt(receiver1); err != nil {
		return err
	}
	wc, err = client.Data()
	if err != nil {
		return err
	}
	_, err = wc.Write(messageBody(sender, receiver1))
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}
	if err = client.Reset(); err != nil {
		return err
	}
	if err = client.Mail(sender); err != nil {
		return err
	}
	if err = client.Rcpt(receiver2); err != nil {
		return err
	}
	if err = client.Rcpt(receiver3); err != nil {
		return err
	}
	wc, err = client.Data()
	if err != nil {
		return err
	}
	_, err = wc.Write(messageBody(sender, receiver2))
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

// Runs successful SMTP session with target host
func runSuccessfulSMTPSession(hostAddress string, portNumber int, fullFlow bool) error {
	connection, _ := net.DialTimeout(networkProtocol, serverWithPortNumber(hostAddress, portNumber), time.Duration(2)*time.Second)
	client, _ := smtp.NewClient(connection, hostAddress)
	var err error

	if err = client.Hello("olo.com"); err != nil {
		return err
	}

	if fullFlow {
		if err = runFullFlow(client); err != nil {
			return err
		}
	}

	if err = client.Quit(); err != nil {
		return err
	}
	if err = client.Close(); err != nil {
		return err
	}

	return nil
}
