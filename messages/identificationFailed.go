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

var identificationFailedMap = make(map[uint]*identificationFailed)

func (f *identificationFailed) GetNOA(instance uint) Message {
	identificationFailed_, ok := identificationFailedMap[instance]

	if ok {
		return identificationFailed_
	}

	identificationFailedMap[instance] = &identificationFailed{PacketId: IdentificationFailedID}
	return identificationFailedMap[instance]
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
