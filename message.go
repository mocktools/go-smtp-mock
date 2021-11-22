package smtpmock

import "sync"

// Structure for storing the result of SMTP client-server interaction
type message struct {
	heloRequest, heloResponse                            string
	mailfromRequest, mailfromResponse                    string
	rcpttoRequest, rcpttoResponse                        string
	dataRequest, dataResponse                            string
	msgRequest, msgResponse                              string
	helo, mailfrom, rcptto, data, msg, cleared, quitSent bool
}

// message methods

// Cleared status predicate. Returns true for case when message struct
// was cleared. Otherwise returns false
func (message *message) isCleared() bool {
	return message.cleared
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
