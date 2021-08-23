package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
)

type trustCertificate struct {
	PacketId uint32
	id       uint32
	hash     string
}

var tCertif = &trustCertificate{PacketId: 2178}

func GetTrustCertificateNOA() *trustCertificate {
	return tCertif
}

func (t *trustCertificate) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, t.id)
	utils.WriteUTF(buff, t.hash)
}

func (t *trustCertificate) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &t.id)
	t.hash = utils.ReadUTF(reader)
}
