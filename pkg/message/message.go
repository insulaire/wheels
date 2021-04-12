package message

type IMessage interface {
	GetLength() uint32
	GetMessageId() uint32
	GetBody() []byte
	SetBody([]byte)
}

type Message struct {
	length    uint32
	messageId uint32
	body      []byte
}

func NewMessage(messageId uint32, body []byte) IMessage {
	return &Message{
		length:    uint32(len(body)),
		messageId: messageId,
		body:      body,
	}
}

func (msg *Message) GetLength() uint32 {
	return msg.length
}
func (msg *Message) GetMessageId() uint32 {
	return msg.messageId
}
func (msg *Message) GetBody() []byte {
	return msg.body
}
func (msg *Message) SetBody(buf []byte) {
	msg.body = buf
}
