package messages

import (
	"bytes"
	"fmt"
)

type identificationFailedForBadVersion struct {
	PacketId uint32
	Idf      *identificationFailed
	Version  *version
}

var identificationFailedForBadVersionMap = make(map[uint]*identificationFailedForBadVersion)

func GetIdentificationFailedForBadVersionNOA(instance uint) *identificationFailedForBadVersion {
	identificationFailedForBadVersion_, ok := identificationFailedForBadVersionMap[instance]

	if ok {
		return identificationFailedForBadVersion_
	}

	identificationFailedForBadVersionMap[instance] = &identificationFailedForBadVersion{PacketId: IdentificationFailedForBadVersionID,
		Idf:     new(identificationFailed),
		Version: new(version)}
	return identificationFailedForBadVersionMap[instance]
}

func (f *identificationFailedForBadVersion) Serialize(buff *bytes.Buffer) {
	f.Idf.Serialize(buff)
	f.Version.Serialize(buff)
}

func (f *identificationFailedForBadVersion) Deserialize(reader *bytes.Reader) {
	f.Idf.Deserialize(reader)
	f.Version.Deserialize(reader)
}

func (f *identificationFailedForBadVersion) GetPacketId() uint32 {
	return f.PacketId
}

func (f *identificationFailedForBadVersion) String() string {
	return fmt.Sprintf("PacketId: %d\nVersion: %v\nidentificationFailed: %v\n", f.PacketId, f.Version, f.Idf)
}
