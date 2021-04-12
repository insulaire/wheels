package message

import (
	"bytes"
	"encoding/binary"
)

type IPack interface {
	Pack(IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
	GetHeaderLength() uint32
}

type Pack struct{}

func NewPack() IPack {
	return &Pack{}
}
func (p *Pack) GetHeaderLength() uint32 {
	return 8
}

func (p *Pack) Pack(msg IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, msg.GetLength()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMessageId()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetBody()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Pack) UnPack(buf []byte) (IMessage, error) {
	reader := bytes.NewReader(buf)
	msg := &Message{}
	if err := binary.Read(reader, binary.LittleEndian, &msg.length); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &msg.messageId); err != nil {
		return nil, err
	}
	return msg, nil
}
