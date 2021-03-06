// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-09-11 12:48:05.4908298 +0200 CEST m=+37.646018701

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type FightStartingPositions struct {
	PacketId                uint32
	positionsForChallengers []int32
	positionsForDefenders   []int32
}

var FightStartingPositionsMap = make(map[uint]*FightStartingPositions)

func (Fi *FightStartingPositions) GetNOA(instance uint) Message {
	FightStartingPositions_, ok := FightStartingPositionsMap[instance]

	if ok {
		return FightStartingPositions_
	}

	FightStartingPositionsMap[instance] = &FightStartingPositions{PacketId: FightStartingPositionsID}
	return FightStartingPositionsMap[instance]
}

func (Fi *FightStartingPositions) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(Fi.positionsForChallengers)))
	for i := 0; i < len(Fi.positionsForChallengers); i++ {
		utils.WriteVarInt16(buff, Fi.positionsForChallengers[i])
	}
	_ = binary.Write(buff, binary.BigEndian, uint16(len(Fi.positionsForDefenders)))
	for i := 0; i < len(Fi.positionsForDefenders); i++ {
		utils.WriteVarInt16(buff, Fi.positionsForDefenders[i])
	}
}

func (Fi *FightStartingPositions) Deserialize(reader *bytes.Reader) {
	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	Fi.positionsForChallengers = nil
	for i := 0; i < int(len0_); i++ {
		Fi.positionsForChallengers = append(Fi.positionsForChallengers, utils.ReadVarInt16(reader))
	}
	var len1_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len1_)
	Fi.positionsForDefenders = nil
	for i := 0; i < int(len1_); i++ {
		Fi.positionsForDefenders = append(Fi.positionsForDefenders, utils.ReadVarInt16(reader))
	}
}

func (Fi *FightStartingPositions) GetPacketId() uint32 {
	return Fi.PacketId
}

func (Fi *FightStartingPositions) String() string {
	return fmt.Sprintf("packetId: %d\n", Fi.PacketId)
}
