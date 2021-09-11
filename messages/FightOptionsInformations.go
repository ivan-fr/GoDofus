// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-09-11 12:31:27.8027686 +0200 CEST m=+31.249836601

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type FightOptionsInformations struct {
	PacketId                uint32
	isSecret                bool
	isRestrictedToPartyOnly bool
	isClosed                bool
	isAskingForHelp         bool
}

var FightOptionsInformationsMap = make(map[uint]*FightOptionsInformations)

func (Fi *FightOptionsInformations) GetNOA(instance uint) Message {
	FightOptionsInformations_, ok := FightOptionsInformationsMap[instance]

	if ok {
		return FightOptionsInformations_
	}

	FightOptionsInformationsMap[instance] = &FightOptionsInformations{PacketId: FightOptionsInformationsID}
	return FightOptionsInformationsMap[instance]
}

func (Fi *FightOptionsInformations) Serialize(buff *bytes.Buffer) {
	var box0 uint32
	box0 = utils.SetFlag(box0, 0, Fi.isSecret)
	box0 = utils.SetFlag(box0, 1, Fi.isRestrictedToPartyOnly)
	box0 = utils.SetFlag(box0, 2, Fi.isClosed)
	box0 = utils.SetFlag(box0, 3, Fi.isAskingForHelp)
	_ = binary.Write(buff, binary.BigEndian, byte(box0))
}

func (Fi *FightOptionsInformations) Deserialize(reader *bytes.Reader) {
	var box0 byte
	_ = binary.Read(reader, binary.BigEndian, &box0)
	Fi.isSecret = utils.GetFlag(uint32(box0), 0)
	Fi.isRestrictedToPartyOnly = utils.GetFlag(uint32(box0), 1)
	Fi.isClosed = utils.GetFlag(uint32(box0), 2)
	Fi.isAskingForHelp = utils.GetFlag(uint32(box0), 3)
}

func (Fi *FightOptionsInformations) GetPacketId() uint32 {
	return Fi.PacketId
}

func (Fi *FightOptionsInformations) String() string {
	return fmt.Sprintf("packetId: %d\n", Fi.PacketId)
}
