package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type Identification struct {
	PacketId            uint32
	Version             *version
	Lang                []byte
	Credentials         []byte
	ServerId            uint16
	AutoSelectServer    bool
	UseCertificate      bool
	UseLoginToken       bool
	SessionOptionalSalt float64
	FailedAttempts      []uint32
	AesKEY_             []byte
}

var identificationMap = make(map[uint]*Identification)

func (id *Identification) GetNOA(instance uint) Message {
	identification_, ok := identificationMap[instance]

	if ok {
		return identification_
	}

	identificationMap[instance] = &Identification{PacketId: IdentificationID, Version: new(version),
		UseCertificate: false, UseLoginToken: false, ServerId: 0}
	return identificationMap[instance]
}

func (id *Identification) Serialize(buff *bytes.Buffer) {
	var box uint32
	box = utils.SetFlag(box, 0, id.AutoSelectServer)
	box = utils.SetFlag(box, 1, id.UseCertificate)
	box = utils.SetFlag(box, 2, id.UseLoginToken)

	_ = binary.Write(buff, binary.BigEndian, byte(box))
	id.Version.Serialize(buff)
	utils.WriteUTF(buff, id.Lang)
	utils.WriteVarInt32(buff, int32(len(id.Credentials)))
	_ = binary.Write(buff, binary.BigEndian, id.Credentials)
	_ = binary.Write(buff, binary.BigEndian, id.ServerId)

	if id.SessionOptionalSalt < -9007199254740992 || id.SessionOptionalSalt > 9007199254740992 {
		panic("Forbidden value on element SessionOptionalSalt.")
	}

	utils.WriteVarInt64(buff, id.SessionOptionalSalt)

	_ = binary.Write(buff, binary.BigEndian, uint16(len(id.FailedAttempts)))

	for i := 0; i < len(id.FailedAttempts); i++ {
		utils.WriteVarInt16(buff, int32(id.FailedAttempts[i]))
	}
}

func (id *Identification) Deserialize(reader *bytes.Reader) {
	var box byte
	_ = binary.Read(reader, binary.BigEndian, &box)
	id.AutoSelectServer = utils.GetFlag(uint32(box), 0)
	id.UseCertificate = utils.GetFlag(uint32(box), 1)
	id.UseLoginToken = utils.GetFlag(uint32(box), 2)
	id.Version.Deserialize(reader)
	id.Lang = utils.ReadUTF(reader)

	credLen := utils.ReadVarInt32(reader)
	id.Credentials = make([]byte, credLen)
	_ = binary.Read(reader, binary.BigEndian, id.Credentials)

	_ = binary.Read(reader, binary.BigEndian, &id.ServerId)

	id.SessionOptionalSalt = float64(utils.ReadVarInt64(reader))

	var failsLen uint16
	_ = binary.Read(reader, binary.BigEndian, &failsLen)

	id.FailedAttempts = nil
	for i := uint16(0); i < failsLen; i++ {
		id.FailedAttempts = append(id.FailedAttempts, uint32(utils.ReadVarInt16(reader)))
	}
}

func (id *Identification) GetPacketId() uint32 {
	return id.PacketId
}

func (id *Identification) String() string {
	return fmt.Sprintf("PacketId: %d\nVersion: %v\nidentification: ...\n\t%s\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n",
		id.PacketId, id.Version, id.Lang, id.Credentials, id.SessionOptionalSalt, id.FailedAttempts, id.ServerId, len(id.Credentials))
}