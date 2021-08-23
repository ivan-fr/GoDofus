package messages

import (
	"bytes"
	"fmt"
)

type identificationFailedForBadVersion struct {
	packetId uint32
	Idf      *identificationFailed
	Version  *version
}

var idfv = &identificationFailedForBadVersion{packetId: VersionID, Idf: new(identificationFailed), Version: new(version)}

func GetIdentificationFailedForBadVersionNOA() *identificationFailedForBadVersion {
	return idfv
}

func (f *identificationFailedForBadVersion) Deserialize(reader *bytes.Reader) {
	f.Idf.Deserialize(reader)
	f.Version.Deserialize(reader)
}

func (f *identificationFailedForBadVersion) String() string {
	return fmt.Sprintf("packetId: %d\nVersion: %v\nidentificationFailed: %v\n", f.packetId, f.Version, f.Idf)
}
