package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageHeloRequest(t *testing.T) {
	t.Run("getter for heloRequest field", func(t *testing.T) {
		message := Message{heloRequest: "some context"}

		assert.Equal(t, message.heloRequest, message.HeloRequest())
	})
}

func TestMessageHeloResponse(t *testing.T) {
	t.Run("getter for heloRequest field", func(t *testing.T) {
		message := Message{heloResponse: "some context"}

		assert.Equal(t, message.heloResponse, message.HeloResponse())
	})
}

func TestMessageHelo(t *testing.T) {
	t.Run("getter for helo field", func(t *testing.T) {
		message := Message{helo: true}

		assert.Equal(t, message.helo, message.Helo())
	})
}

func TestMessageMailfromRequest(t *testing.T) {
	t.Run("getter for mailfromRequest field", func(t *testing.T) {
		message := Message{mailfromRequest: "some context"}

		assert.Equal(t, message.mailfromRequest, message.MailfromRequest())
	})
}

func TestMessageMailfromResponse(t *testing.T) {
	t.Run("getter for mailfromResponse field", func(t *testing.T) {
		message := Message{mailfromResponse: "some context"}

		assert.Equal(t, message.mailfromResponse, message.MailfromResponse())
	})
}

func TestMessageMailfrom(t *testing.T) {
	t.Run("getter for mailfrom field", func(t *testing.T) {
		message := Message{mailfrom: true}

		assert.Equal(t, message.mailfrom, message.Mailfrom())
	})
}

func TestMessageRcpttoRequestResponse(t *testing.T) {
	t.Run("getter for rcpttoRequestResponse field", func(t *testing.T) {
		message := Message{rcpttoRequestResponse: [][]string{{"request", "response"}}}

		assert.Equal(t, message.rcpttoRequestResponse, message.RcpttoRequestResponse())
	})
}

func TestMessageRcptto(t *testing.T) {
	t.Run("getter for rcptto field", func(t *testing.T) {
		message := Message{rcptto: true}

		assert.Equal(t, message.rcptto, message.Rcptto())
	})
}

func TestMessageDataRequest(t *testing.T) {
	t.Run("getter for dataRequest field", func(t *testing.T) {
		message := Message{dataRequest: "some context"}

		assert.Equal(t, message.dataRequest, message.DataRequest())
	})
}

func TestMessageDataResponse(t *testing.T) {
	t.Run("getter for dataResponse field", func(t *testing.T) {
		message := Message{dataResponse: "some context"}

		assert.Equal(t, message.dataResponse, message.DataResponse())
	})
}

func TestMessageData(t *testing.T) {
	t.Run("getter for data field", func(t *testing.T) {
		message := Message{data: true}

		assert.Equal(t, message.data, message.Data())
	})
}

func TestMessageMsgRequest(t *testing.T) {
	t.Run("getter for msgRequest field", func(t *testing.T) {
		message := Message{msgRequest: "some context"}

		assert.Equal(t, message.msgRequest, message.MsgRequest())
	})
}

func TestMessageMsgResponse(t *testing.T) {
	t.Run("getter for msgRequest field", func(t *testing.T) {
		message := Message{msgResponse: "some context"}

		assert.Equal(t, message.msgResponse, message.MsgResponse())
	})
}

func TestMessageMsg(t *testing.T) {
	t.Run("getter for msg field", func(t *testing.T) {
		message := Message{msg: true}

		assert.Equal(t, message.msg, message.Msg())
	})
}

func TestMessageRsetRequest(t *testing.T) {
	t.Run("getter for rsetRequest field", func(t *testing.T) {
		message := Message{rsetRequest: "some context"}

		assert.Equal(t, message.rsetRequest, message.RsetRequest())
	})
}

func TestMessageRsetResponse(t *testing.T) {
	t.Run("getter for rsetRequest field", func(t *testing.T) {
		message := Message{rsetResponse: "some context"}

		assert.Equal(t, message.rsetResponse, message.RsetResponse())
	})
}

func TestMessageRset(t *testing.T) {
	t.Run("getter for rset field", func(t *testing.T) {
		message := Message{rset: true}

		assert.Equal(t, message.rset, message.Rset())
	})
}

func TestMessageNoop(t *testing.T) {
	t.Run("getter for noop field", func(t *testing.T) {
		message := Message{noop: true}

		assert.Equal(t, message.noop, message.Noop())
	})
}

func TestMessageQuitSent(t *testing.T) {
	t.Run("getter for quitSent field", func(t *testing.T) {
		message := Message{quitSent: true}

		assert.Equal(t, message.quitSent, message.QuitSent())
	})
}

func TestMessageIsConsistent(t *testing.T) {
	t.Run("when consistent", func(t *testing.T) {
		message := &Message{mailfrom: true, rcptto: true, data: true, msg: true}

		assert.True(t, message.IsConsistent())
	})

	t.Run("when not consistent MAILFROM", func(t *testing.T) {

		assert.False(t, new(Message).IsConsistent())
	})

	t.Run("when not consistent RCPTTO", func(t *testing.T) {
		message := &Message{mailfrom: true}

		assert.False(t, message.IsConsistent())
	})

	t.Run("when not consistent DATA", func(t *testing.T) {
		message := &Message{mailfrom: true, rcptto: true}

		assert.False(t, message.IsConsistent())
	})

	t.Run("when not consistent MSG", func(t *testing.T) {
		message := &Message{mailfrom: true, rcptto: true, data: true}

		assert.False(t, message.IsConsistent())
	})
}

func TestMessagePointerIsConsistent(t *testing.T) {
	t.Run("when consistent", func(t *testing.T) {
		message := &Message{mailfrom: true, rcptto: true, data: true, msg: true}

		assert.True(t, message.IsConsistent())
	})

	t.Run("when not consistent MAILFROM", func(t *testing.T) {

		assert.False(t, new(Message).IsConsistent())
	})

	t.Run("when not consistent RCPTTO", func(t *testing.T) {
		message := &Message{mailfrom: true}

		assert.False(t, message.IsConsistent())
	})

	t.Run("when not consistent DATA", func(t *testing.T) {
		message := &Message{mailfrom: true, rcptto: true}

		assert.False(t, message.IsConsistent())
	})

	t.Run("when not consistent MSG", func(t *testing.T) {
		message := &Message{mailfrom: true, rcptto: true, data: true}

		assert.False(t, message.IsConsistent())
	})
}

func TestMessageIsIncludesSuccessfulRcpttoResponse(t *testing.T) {
	targetSuccessfulResponse := "response"

	t.Run("when successful RCPTTO response exists", func(t *testing.T) {
		message := &Message{rcpttoRequestResponse: [][]string{{"request", targetSuccessfulResponse}}}

		assert.True(t, message.isIncludesSuccessfulRcpttoResponse(targetSuccessfulResponse))
	})

	t.Run("when successful RCPTTO response not exists", func(t *testing.T) {
		assert.False(t, new(Message).isIncludesSuccessfulRcpttoResponse(targetSuccessfulResponse))
	})
}

func TestMessagesAppend(t *testing.T) {
	t.Run("addes message pointer into items slice", func(t *testing.T) {
		message, messages := new(Message), new(messages)
		messages.append(message)

		messages.RLock()
		assert.Same(t, message, messages.items[0])
		messages.RUnlock()
	})
}

func TestMessagesCopy(t *testing.T) {
	t.Run("copies messages", func(t *testing.T) {
		message, messages := new(Message), new(messages)
		message.heloRequest = "foobar"
		messages.append(message)
		copyMessages := messages.copy()

		assert.Len(t, copyMessages, 1)
		assert.Equal(t, *message, copyMessages[0])
	})
}
