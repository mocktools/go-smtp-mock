package smtpmock

import "github.com/alecthomas/repr"

func (server *Server) Count() int {
	return len(server.messages.items)
}

func (server *Server) Message(index int) string {
	if server.Count() < index {
		return ""
	}
	return repr.String(server.messages.items[index])
}
