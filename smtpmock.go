package smtpmock

// Builds new SMTP mock server, based on passed configuration attributes
func New(config ConfigurationAttr) *server {
	return newServer(NewConfiguration(config))
}
