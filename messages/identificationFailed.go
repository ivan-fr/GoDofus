package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type identificationFailed struct {
	packetId uint32
	Reason   uint8
}

var idf = &identificationFailed{packetId: IdentificationFailedID}

func GetIdentificationFailedNOA() *identificationFailed {
	return idf
}

func (f *identificationFailed) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &f.Reason)
}

func (f *identificationFailed) String(reader *bytes.Reader) string {
	return fmt.Sprintf("REASEON %d\n", f.Reason)
}
