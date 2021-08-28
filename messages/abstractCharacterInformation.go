// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 21:13:24.9871213 +0200 CEST m=+0.019936701

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type abstractCharacterInformation struct {
	PacketId uint32
	id       float64
}

var abstractCharacterInformationMap = make(map[uint]*abstractCharacterInformation)

func GetAbstractCharacterInformationNOA(instance uint) *abstractCharacterInformation {
	abstractCharacterInformation_, ok := abstractCharacterInformationMap[instance]

	if ok {
		return abstractCharacterInformation_
	}

	abstractCharacterInformationMap[instance] = &abstractCharacterInformation{PacketId: AbstractCharacterInformationID}
	return abstractCharacterInformationMap[instance]
}

func (a *abstractCharacterInformation) Serialize(buff *bytes.Buffer) {
	utils.WriteVarLong(buff, a.id)
}

func (a *abstractCharacterInformation) Deserialize(reader *bytes.Reader) {
	a.id = float64(utils.ReadVarUInt64(reader))
}

func (a *abstractCharacterInformation) GetPacketId() uint32 {
	return a.PacketId
}

func (a *abstractCharacterInformation) String() string {
	return fmt.Sprintf("packetId: %d\n", a.PacketId)
}
