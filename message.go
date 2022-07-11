package smtpmock

import "sync"

// Structure for storing the result of SMTP client-server interaction. Context-included
// commands should be represented as request/response structure fields
type message struct {
	heloRequest, heloResponse                         string
	mailfromRequest, mailfromResponse                 string
	rcpttoRequest, rcpttoResponse                     string
	dataRequest, dataResponse                         string
	msgRequest, msgResponse                           string
	rsetRequest, rsetResponse                         string
	helo, mailfrom, rcptto, data, msg, rset, quitSent bool
}

// message methods

// message getters

// Getter for heloRequest field
func (message *message) HeloRequest() string {
	return message.heloRequest
}

// Getter for heloResponse field
func (message *message) HeloResponse() string {
	return message.heloResponse
}

// Getter for helo field
func (message *message) Helo() bool {
	return message.helo
}

// Getter for mailfromRequest field
func (message *message) MailfromRequest() string {
	return message.mailfromRequest
}

// Getter for mailfromResponse field
func (message *message) MailfromResponse() string {
	return message.mailfromResponse
}

// Getter for mailfrom field
func (message *message) Mailfrom() bool {
	return message.mailfrom
}

// Getter for rcpttoRequest field
func (message *message) RcpttoRequest() string {
	return message.rcpttoRequest
}

// Getter for rcpttoResponse field
func (message *message) RcpttoResponse() string {
	return message.rcpttoResponse
}

// Getter for rcptto field
func (message *message) Rcptto() bool {
	return message.rcptto
}

// Getter for dataRequest field
func (message *message) DataRequest() string {
	return message.dataRequest
}

// Getter for dataResponse field
func (message *message) DataResponse() string {
	return message.dataResponse
}

// Getter for data field
func (message *message) Data() bool {
	return message.data
}

// Getter for msgRequest field
func (message *message) MsgRequest() string {
	return message.msgRequest
}

// Getter for msgResponse field
func (message *message) MsgResponse() string {
	return message.msgResponse
}

// Getter for msg field
func (message *message) Msg() bool {
	return message.msg
}

// Getter for rsetRequest field
func (message *message) RsetRequest() string {
	return message.rsetRequest
}

// Getter for rsetResponse field
func (message *message) RsetResponse() string {
	return message.rsetResponse
}

// Getter for rset field
func (message *message) Rset() bool {
	return message.rset
}

// Getter for quitSent field
func (message *message) QuitSent() bool {
	return message.quitSent
}

// Message consistency status predicate. Returns true for case when message struct is consistent.
// It means that MAILFROM, RCPTTO, DATA commands and message context were successful.
// Otherwise returns false
func (message *message) isConsistent() bool {
	return message.mailfrom && message.rcptto && message.data && message.msg
}

// Pointer to empty message
var zeroMessage = &message{}

// Concurrent type that can be safely shared between goroutines
type messages struct {
	sync.RWMutex
	items []*message
}

// messages methods

// Addes new message pointer into concurrent messages slice
func (messages *messages) append(item *message) {
	messages.Lock()
	defer messages.Unlock()

	messages.items = append(messages.items, item)
}
