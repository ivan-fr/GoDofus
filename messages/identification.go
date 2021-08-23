package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
)

type identification struct {
	PacketId            uint32
	version             version
	lang                string
	credentials         []byte
	serverId            uint16
	autoConnect         bool
	useCertificate      bool
	useLoginToken       bool
	sessionOptionalSalt float64
	failedAttempts      []uint32
}

var idCation = &identification{PacketId: 2767}

func GetIdentificationNOA() *identification {
	return idCation
}

func (id *identification) Serialize(buff *bytes.Buffer) {
	var box uint32
	box = utils.SetFlag(box, 1, id.autoConnect)
	box = utils.SetFlag(box, 2, id.useCertificate)
	box = utils.SetFlag(box, 3, id.useLoginToken)

	_ = binary.Write(buff, binary.BigEndian, byte(box))
	id.version.Serialize(buff)
	utils.WriteUTF(buff, id.lang)
	utils.WriteVarInt32(buff, int32(len(id.credentials)))
	_ = binary.Write(buff, binary.BigEndian, id.credentials)
	_ = binary.Write(buff, binary.BigEndian, id.serverId)

	if id.sessionOptionalSalt < -9007199254740992 || id.sessionOptionalSalt > 9007199254740992 {
		panic("Forbidden value on element sessionOptionalSalt.")
	}

}
