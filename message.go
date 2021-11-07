package smtpmock

import "sync"

// Structure for storing the result of SMTP client-server interaction
type message struct {
	heloRequest, heloResponse         string
	mailfromRequest, mailfromResponse string
	rcpttoRequest, rcpttoResponse     string
	dataRequest, dataResponse         string
	msgRequest, msgResponse           string
	helo, mailfrom, rcptto, data, msg bool
}

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
