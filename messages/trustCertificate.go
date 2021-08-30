package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
)

type trustCertificate struct {
	PacketId uint32
	id       uint32
	hash     []byte
}

var trustCertificateMap = make(map[uint]*trustCertificate)

func (t *trustCertificate) GetNOA(instance uint) Message {
	trustCertificate_, ok := trustCertificateMap[instance]

	if ok {
		return trustCertificate_
	}

	trustCertificateMap[instance] = &trustCertificate{PacketId: TrustCertificateID}
	return trustCertificateMap[instance]
}

func (t *trustCertificate) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, t.id)
	utils.WriteUTF(buff, t.hash)
}

func (t *trustCertificate) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &t.id)
	t.hash = utils.ReadUTF(reader)
}

func (t *trustCertificate) GetPacketId() uint32 {
	return t.PacketId
}
