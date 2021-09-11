// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-09-11 12:34:40.8027408 +0200 CEST m=+20.981319201

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type FightTeamMemberInformations struct {
	PacketId uint32
	id       float64
}

var FightTeamMemberInformationsMap = make(map[uint]*FightTeamMemberInformations)

func (Fi *FightTeamMemberInformations) GetNOA(instance uint) Message {
	FightTeamMemberInformations_, ok := FightTeamMemberInformationsMap[instance]

	if ok {
		return FightTeamMemberInformations_
	}

	FightTeamMemberInformationsMap[instance] = &FightTeamMemberInformations{PacketId: FightTeamMemberInformationsID}
	return FightTeamMemberInformationsMap[instance]
}

func (Fi *FightTeamMemberInformations) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, Fi.id)
}

func (Fi *FightTeamMemberInformations) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &Fi.id)
}

func (Fi *FightTeamMemberInformations) GetPacketId() uint32 {
	return Fi.PacketId
}

func (Fi *FightTeamMemberInformations) String() string {
	return fmt.Sprintf("packetId: %d\n", Fi.PacketId)
}
