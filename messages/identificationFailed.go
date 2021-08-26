package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type identificationFailed struct {
	PacketId uint32
	Reason   uint8
}

var idf = &identificationFailed{PacketId: IdentificationFailedID}

func GetIdentificationFailedNOA() *identificationFailed {
	return idf
}

func (f *identificationFailed) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, f.Reason)
}

func (f *identificationFailed) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &f.Reason)
}

func (f *identificationFailed) GetPacketId() uint32 {
	return f.PacketId
}

func (f *identificationFailed) String() string {
	return fmt.Sprintf("Reason %d\n", f.Reason)
}
