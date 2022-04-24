package smtpmock

import "github.com/alecthomas/repr"

func (server *Server) Count() int {
	return len(server.messages.items)
}

func (server *Server) Message(index uint) string {
	if int(index) < server.Count() {
		return repr.String(server.messages.items[index])
	}

	return ""
}

func (server *Server) Clear() {
	server.messages.items = nil
}
