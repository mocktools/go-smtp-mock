package smtpmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageHeloRequest(t *testing.T) {
	t.Run("getter for heloRequest field", func(t *testing.T) {
		message := &message{heloRequest: "some context"}

		assert.Equal(t, message.heloRequest, message.HeloRequest())
	})
}

func TestMessageHeloResponse(t *testing.T) {
	t.Run("getter for heloRequest field", func(t *testing.T) {
		message := &message{heloResponse: "some context"}

		assert.Equal(t, message.heloResponse, message.HeloResponse())
	})
}

func TestMessageHelo(t *testing.T) {
	t.Run("getter for helo field", func(t *testing.T) {
		message := &message{helo: true}

		assert.Equal(t, message.helo, message.Helo())
	})
}

func TestMessageMailfromRequest(t *testing.T) {
	t.Run("getter for mailfromRequest field", func(t *testing.T) {
		message := &message{mailfromRequest: "some context"}

		assert.Equal(t, message.mailfromRequest, message.MailfromRequest())
	})
}

func TestMessageMailfromResponse(t *testing.T) {
	t.Run("getter for mailfromResponse field", func(t *testing.T) {
		message := &message{mailfromResponse: "some context"}

		assert.Equal(t, message.mailfromResponse, message.MailfromResponse())
	})
}

func TestMessageMailfrom(t *testing.T) {
	t.Run("getter for mailfrom field", func(t *testing.T) {
		message := &message{mailfrom: true}

		assert.Equal(t, message.mailfrom, message.Mailfrom())
	})
}

func TestMessageRcpttoRequest(t *testing.T) {
	t.Run("getter for rcpttoRequest field", func(t *testing.T) {
		message := &message{rcpttoRequest: "some context"}

		assert.Equal(t, message.rcpttoRequest, message.RcpttoRequest())
	})
}

func TestMessageRcpttoResponse(t *testing.T) {
	t.Run("getter for rcpttoResponse field", func(t *testing.T) {
		message := &message{rcpttoResponse: "some context"}

		assert.Equal(t, message.rcpttoResponse, message.RcpttoResponse())
	})
}

func TestMessageRcptto(t *testing.T) {
	t.Run("getter for rcptto field", func(t *testing.T) {
		message := &message{rcptto: true}

		assert.Equal(t, message.rcptto, message.Rcptto())
	})
}

func TestMessageDataRequest(t *testing.T) {
	t.Run("getter for dataRequest field", func(t *testing.T) {
		message := &message{dataRequest: "some context"}

		assert.Equal(t, message.dataRequest, message.DataRequest())
	})
}

func TestMessageDataResponse(t *testing.T) {
	t.Run("getter for dataResponse field", func(t *testing.T) {
		message := &message{dataResponse: "some context"}

		assert.Equal(t, message.dataResponse, message.DataResponse())
	})
}

func TestMessageData(t *testing.T) {
	t.Run("getter for data field", func(t *testing.T) {
		message := &message{data: true}

		assert.Equal(t, message.data, message.Data())
	})
}

func TestMessageMsgRequest(t *testing.T) {
	t.Run("getter for msgRequest field", func(t *testing.T) {
		message := &message{msgRequest: "some context"}

		assert.Equal(t, message.msgRequest, message.MsgRequest())
	})
}

func TestMessageMsgResponse(t *testing.T) {
	t.Run("getter for msgRequest field", func(t *testing.T) {
		message := &message{msgResponse: "some context"}

		assert.Equal(t, message.msgResponse, message.MsgResponse())
	})
}

func TestMessageMsg(t *testing.T) {
	t.Run("getter for msg field", func(t *testing.T) {
		message := &message{msg: true}

		assert.Equal(t, message.msg, message.Msg())
	})
}

func TestMessageRsetRequest(t *testing.T) {
	t.Run("getter for rsetRequest field", func(t *testing.T) {
		message := &message{rsetRequest: "some context"}

		assert.Equal(t, message.rsetRequest, message.RsetRequest())
	})
}

func TestMessageRsetResponse(t *testing.T) {
	t.Run("getter for rsetRequest field", func(t *testing.T) {
		message := &message{rsetResponse: "some context"}

		assert.Equal(t, message.rsetResponse, message.RsetResponse())
	})
}

func TestMessageRset(t *testing.T) {
	t.Run("getter for rset field", func(t *testing.T) {
		message := &message{rset: true}

		assert.Equal(t, message.rset, message.Rset())
	})
}

func TestMessageQuitSent(t *testing.T) {
	t.Run("getter for quitSent field", func(t *testing.T) {
		message := &message{quitSent: true}

		assert.Equal(t, message.quitSent, message.QuitSent())
	})
}

func TestMessageIsConsistent(t *testing.T) {
	t.Run("when consistent", func(t *testing.T) {
		message := &message{mailfrom: true, rcptto: true, data: true, msg: true}

		assert.True(t, message.isConsistent())
	})

	t.Run("when not consistent MAILFROM", func(t *testing.T) {

		assert.False(t, new(message).isConsistent())
	})

	t.Run("when not consistent RCPTTO", func(t *testing.T) {
		message := &message{mailfrom: true}

		assert.False(t, message.isConsistent())
	})

	t.Run("when not consistent DATA", func(t *testing.T) {
		message := &message{mailfrom: true, rcptto: true}

		assert.False(t, message.isConsistent())
	})

	t.Run("when not consistent MSG", func(t *testing.T) {
		message := &message{mailfrom: true, rcptto: true, data: true}

		assert.False(t, message.isConsistent())
	})
}

func TestMessagesAppend(t *testing.T) {
	t.Run("addes message pointer into items slice", func(t *testing.T) {
		message, messages := new(message), new(messages)
		messages.append(message)

		assert.Same(t, message, messages.items[0])
	})
}
