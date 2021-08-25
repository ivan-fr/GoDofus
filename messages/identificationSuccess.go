// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-25 10:00:15.98031 +0200 CEST m=+0.002609801

package messages

import (
	"GoDofus/structs"
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type identificationSuccess struct {
	PacketId                    uint32
	hasRights                   bool
	hasConsoleRight             bool
	wasAlreadyConnected         bool
	login                       []byte
	aTI                         *structs.AccountTagInformation
	accoutId                    uint32
	communityId                 byte
	secretQuestion              []byte
	accountCreation             float64
	subscriptionElapsedDuration float64
	subscriptionEndDate         float64
	havenbagAvailableRoom       uint8
}

var identificationSuccess_ = &identificationSuccess{PacketId: IdentificationSuccessID}

func GetIdentificationSuccessNOA() *identificationSuccess {
	return identificationSuccess_
}

func (i *identificationSuccess) Serialize(buff *bytes.Buffer) {

}

func (i *identificationSuccess) Deserialize(reader *bytes.Reader) {
	var box byte
	_ = binary.Read(reader, binary.BigEndian, &box)

	i.hasRights = utils.GetFlag(uint32(box), 0)
	i.hasConsoleRight = utils.GetFlag(uint32(box), 1)
	i.wasAlreadyConnected = utils.GetFlag(uint32(box), 2)

	i.login = utils.ReadUTF(reader)

	i.aTI = new(structs.AccountTagInformation)
	i.aTI.Deserialize(reader)

	_ = binary.Read(reader, binary.BigEndian, &i.accoutId)
	_ = binary.Read(reader, binary.BigEndian, &i.communityId)

	i.secretQuestion = utils.ReadUTF(reader)

	_ = binary.Read(reader, binary.BigEndian, &i.accountCreation)
	_ = binary.Read(reader, binary.BigEndian, &i.subscriptionElapsedDuration)
	_ = binary.Read(reader, binary.BigEndian, &i.subscriptionEndDate)
	_ = binary.Read(reader, binary.BigEndian, &i.havenbagAvailableRoom)
}

func (i *identificationSuccess) GetPacketId() uint32 {
	return i.PacketId
}

func (i *identificationSuccess) String(reader *bytes.Reader) string {
	return fmt.Sprintf("PacketId: %d\n", i.PacketId)
}