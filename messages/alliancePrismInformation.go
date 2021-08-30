// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 09:23:54.081485 +0200 CEST m=+9.583237301

package messages

import (
	"bytes"
	"fmt"
)

type alliancePrismInformation struct {
	PacketId            uint32
	prismInformation    *prismInformation
	allianceInformation *allianceInformation
}

var alliancePrismInformationMap = make(map[uint]*alliancePrismInformation)

func GetAlliancePrismInformationNOA(instance uint) *alliancePrismInformation {
	alliancePrismInformation_, ok := alliancePrismInformationMap[instance]

	if ok {
		return alliancePrismInformation_
	}

	alliancePrismInformationMap[instance] = &alliancePrismInformation{PacketId: AlliancePrismInformationID}
	return alliancePrismInformationMap[instance]
}

func (al *alliancePrismInformation) Serialize(buff *bytes.Buffer) {
	al.prismInformation.Serialize(buff)
	al.allianceInformation.Serialize(buff)
}

func (al *alliancePrismInformation) Deserialize(reader *bytes.Reader) {
	al.prismInformation = new(prismInformation)
	al.prismInformation.Deserialize(reader)
	al.allianceInformation = new(allianceInformation)
	al.allianceInformation.Deserialize(reader)
}

func (al *alliancePrismInformation) GetPacketId() uint32 {
	return al.PacketId
}

func (al *alliancePrismInformation) String() string {
	return fmt.Sprintf("packetId: %d\n", al.PacketId)
}