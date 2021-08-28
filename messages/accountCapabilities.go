// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 17:18:46.3839133 +0200 CEST m=+0.020003901

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type accountCapabilities struct {
	PacketId              uint32
	tutorialAvailable     bool
	canCreateNewCharacter bool
	accountId             uint32
	breedsVisible         uint32
	breedsAvailable       int32
	status                byte
}

var accountCapabilitiesMap = make(map[uint]*accountCapabilities)

func GetAccountCapabilitiesNOA(instance uint) *accountCapabilities {
	accountCapabilities_, ok := accountCapabilitiesMap[instance]

	if ok {
		return accountCapabilities_
	}

	accountCapabilitiesMap[instance] = &accountCapabilities{PacketId: AccountCapabilitiesID}
	return accountCapabilitiesMap[instance]
}

func (a *accountCapabilities) Serialize(buff *bytes.Buffer) {
	var box uint32
	box = utils.SetFlag(box, 0, a.tutorialAvailable)
	box = utils.SetFlag(box, 1, a.canCreateNewCharacter)

	_ = binary.Write(buff, binary.BigEndian, byte(box))
	_ = binary.Write(buff, binary.BigEndian, a.accountId)
	_ = binary.Write(buff, binary.BigEndian, a.breedsVisible)

	utils.WriteVarInt32(buff, a.breedsAvailable)
	_ = binary.Write(buff, binary.BigEndian, a.status)
}

func (a *accountCapabilities) Deserialize(reader *bytes.Reader) {
	var box byte
	_ = binary.Read(reader, binary.BigEndian, &box)

	a.tutorialAvailable = utils.GetFlag(uint32(box), 0)
	a.canCreateNewCharacter = utils.GetFlag(uint32(box), 1)

	_ = binary.Read(reader, binary.BigEndian, &a.accountId)
	_ = binary.Read(reader, binary.BigEndian, &a.breedsVisible)

	a.breedsAvailable = utils.ReadVarInt32(reader)
	_ = binary.Read(reader, binary.BigEndian, &a.status)
}

func (a *accountCapabilities) GetPacketId() uint32 {
	return a.PacketId
}

func (a *accountCapabilities) String() string {
	return fmt.Sprintf("packetId: %d\n", a.PacketId)
}
