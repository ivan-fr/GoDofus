package messages

import (
	"bytes"
	"encoding/binary"
)

type version struct {
	packetId  uint32
	Major     uint8
	Minor     uint8
	Code      uint8
	Build     uint32
	BuildType uint8
}

var ve = &version{packetId: VersionID}

func GetVersionNOA() *version {
	return ve
}

func (v *version) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, v.Major)
	_ = binary.Write(buff, binary.BigEndian, v.Minor)
	_ = binary.Write(buff, binary.BigEndian, v.Code)
	_ = binary.Write(buff, binary.BigEndian, v.Build)
	_ = binary.Write(buff, binary.BigEndian, v.BuildType)
}

func (v *version) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &v.Major)
	_ = binary.Read(reader, binary.BigEndian, &v.Minor)
	_ = binary.Read(reader, binary.BigEndian, &v.Code)
	_ = binary.Read(reader, binary.BigEndian, &v.Build)
	_ = binary.Read(reader, binary.BigEndian, &v.BuildType)
}
