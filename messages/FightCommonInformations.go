// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-09-11 12:43:22.1948343 +0200 CEST m=+107.827733001

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type FightCommonInformations struct {
	PacketId                  uint32
	fightId                   int32
	fightType                 byte
	FightTeamInformations2    []*FightTeamInformations
	fightTeamsPositions       []int32
	FightOptionsInformations4 []*FightOptionsInformations
}

var FightCommonInformationsMap = make(map[uint]*FightCommonInformations)

func (Fi *FightCommonInformations) GetNOA(instance uint) Message {
	FightCommonInformations_, ok := FightCommonInformationsMap[instance]

	if ok {
		return FightCommonInformations_
	}

	FightCommonInformationsMap[instance] = &FightCommonInformations{PacketId: FightCommonInformationsID}
	return FightCommonInformationsMap[instance]
}

func (Fi *FightCommonInformations) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt16(buff, Fi.fightId)
	_ = binary.Write(buff, binary.BigEndian, Fi.fightType)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(Fi.FightTeamInformations2)))
	for i := 0; i < len(Fi.FightTeamInformations2); i++ {
		Fi.FightTeamInformations2[i].Serialize(buff)
	}
	_ = binary.Write(buff, binary.BigEndian, uint16(len(Fi.fightTeamsPositions)))
	for i := 0; i < len(Fi.fightTeamsPositions); i++ {
		utils.WriteVarInt16(buff, Fi.fightTeamsPositions[i])
	}
	_ = binary.Write(buff, binary.BigEndian, uint16(len(Fi.FightOptionsInformations4)))
	for i := 0; i < len(Fi.FightOptionsInformations4); i++ {
		Fi.FightOptionsInformations4[i].Serialize(buff)
	}
}

func (Fi *FightCommonInformations) Deserialize(reader *bytes.Reader) {
	Fi.fightId = utils.ReadVarInt16(reader)
	_ = binary.Read(reader, binary.BigEndian, &Fi.fightType)
	var len2_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len2_)
	Fi.FightTeamInformations2 = nil
	for i := 0; i < int(len2_); i++ {
		aMessage2 := new(FightTeamInformations)
		aMessage2.Deserialize(reader)
		Fi.FightTeamInformations2 = append(Fi.FightTeamInformations2, aMessage2)
	}
	var len3_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len3_)
	Fi.fightTeamsPositions = nil
	for i := 0; i < int(len3_); i++ {
		Fi.fightTeamsPositions = append(Fi.fightTeamsPositions, utils.ReadVarInt16(reader))
	}
	var len4_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len4_)
	Fi.FightOptionsInformations4 = nil
	for i := 0; i < int(len4_); i++ {
		aMessage4 := new(FightOptionsInformations)
		aMessage4.Deserialize(reader)
		Fi.FightOptionsInformations4 = append(Fi.FightOptionsInformations4, aMessage4)
	}
}

func (Fi *FightCommonInformations) GetPacketId() uint32 {
	return Fi.PacketId
}

func (Fi *FightCommonInformations) String() string {
	return fmt.Sprintf("packetId: %d\n", Fi.PacketId)
}
