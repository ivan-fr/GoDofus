package messages

import (
	"bytes"
)

type Message interface {
	GetPacketId() uint32
	Serialize(*bytes.Buffer)
	Deserialize(*bytes.Reader)
	GetNOA(uint) Message
}

type Protocol interface {
	GetPacketId() uint32
	Serialize(*bytes.Buffer)
	Deserialize(*bytes.Reader)
}
