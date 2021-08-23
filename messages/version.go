package messages

import (
	"bytes"
	"encoding/binary"
)

type version struct {
	PacketId  uint32
	major     uint8
	minor     uint8
	code      uint8
	build     uint32
	buildType uint8
}

var ve = &version{PacketId: 9475}

func GetVersionNOA() *version {
	return ve
}

func (v *version) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, v.major)
	_ = binary.Write(buff, binary.BigEndian, v.minor)
	_ = binary.Write(buff, binary.BigEndian, v.code)
	_ = binary.Write(buff, binary.BigEndian, v.build)
	_ = binary.Write(buff, binary.BigEndian, v.buildType)
}

func (v *version) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &v.major)
	_ = binary.Read(reader, binary.BigEndian, &v.minor)
	_ = binary.Read(reader, binary.BigEndian, &v.code)
	_ = binary.Read(reader, binary.BigEndian, &v.build)
	_ = binary.Read(reader, binary.BigEndian, &v.buildType)
}
