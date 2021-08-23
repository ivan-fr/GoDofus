package messages

import (
	"bytes"
	"encoding/binary"
)

type identificationFailed struct {
	packetId uint32
	Reason   uint8
}

var idf = &identificationFailed{packetId: VersionID}

func GetIdentificationFailedNOA() *identificationFailed {
	return idf
}

func (f *identificationFailed) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &f.Reason)
}
