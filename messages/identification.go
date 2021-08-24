package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type identification struct {
	packetId            uint32
	Version             *version
	Lang                []byte
	Credentials         []byte
	ServerId            uint16
	AutoSelectServer    bool
	UseCertificate      bool
	UseLoginToken       bool
	SessionOptionalSalt float64
	FailedAttempts      []uint32
}

var idCation = &identification{packetId: IdentificationID, Version: new(version),
	UseCertificate: false, UseLoginToken: false, ServerId: 0}

func GetIdentificationNOA() *identification {
	return idCation
}

func (id *identification) Serialize(buff *bytes.Buffer) {
	var box uint32
	box = utils.SetFlag(box, 1, id.AutoSelectServer)
	box = utils.SetFlag(box, 2, id.UseCertificate)
	box = utils.SetFlag(box, 3, id.UseLoginToken)

	_ = binary.Write(buff, binary.BigEndian, byte(box))
	id.Version.Serialize(buff)
	utils.WriteUTF(buff, id.Lang)
	utils.WriteVarInt32(buff, int32(len(id.Credentials)))
	_ = binary.Write(buff, binary.BigEndian, id.Credentials)
	_ = binary.Write(buff, binary.BigEndian, id.ServerId)

	if id.SessionOptionalSalt < -9007199254740992 || id.SessionOptionalSalt > 9007199254740992 {
		panic("Forbidden value on element SessionOptionalSalt.")
	}

	utils.WriteVarLong(buff, id.SessionOptionalSalt)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(id.FailedAttempts)))

	for i := 0; i < len(id.FailedAttempts); i++ {
		utils.WriteVarShort(buff, int32(id.FailedAttempts[i]))
	}
}

func (id *identification) Deserialize(reader *bytes.Reader) {
	var box byte
	_ = binary.Read(reader, binary.BigEndian, &box)
	id.AutoSelectServer = utils.GetFlag(uint32(box), 0)
	id.UseCertificate = utils.GetFlag(uint32(box), 1)
	id.UseLoginToken = utils.GetFlag(uint32(box), 2)
	id.Version.Deserialize(reader)
	id.Lang = utils.ReadUTF(reader)

	credLen := utils.ReadVarInt32(reader)
	id.Credentials = make([]byte, credLen)
	_ = binary.Read(reader, binary.BigEndian, &id.Credentials)

	_ = binary.Read(reader, binary.BigEndian, &id.ServerId)

	var saltSession int64
	_ = binary.Read(reader, binary.BigEndian, &saltSession)
	id.SessionOptionalSalt = float64(saltSession)

	var failsLen uint16
	_ = binary.Read(reader, binary.BigEndian, &failsLen)
	id.FailedAttempts = nil

	for i := uint16(0); i < failsLen; i++ {
		id.FailedAttempts = append(id.FailedAttempts, uint32(utils.ReadVarInt16(reader)))
	}
}

func (id *identification) String() string {
	return fmt.Sprintf("packetId: %d\nVersion: %v\nidentification: ...\n\t%s\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n",
		id.packetId, id.Version, id.Lang, id.Credentials, id.SessionOptionalSalt, id.FailedAttempts, id.ServerId, len(id.Credentials))
}
