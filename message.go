package smtpmock

// Structure for storing the result of SMTP client-server interaction
type message struct {
	heloRequest, heloResponse string
	helo                      bool
	// mailfromRequest, mailfromResponse string
	// rcpttoRequest, rcpttoResponse     string
	// mailfrom, rcptto bool
}
