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

var idfv = &identificationFailedForBadVersion{PacketId: VersionID, Idf: new(identificationFailed), Version: new(version)}

func GetIdentificationFailedForBadVersionNOA() *identificationFailedForBadVersion {
	return idfv
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
