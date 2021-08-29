// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 22:07:49.8968201 +0200 CEST m=+30.550485701

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type simpleCharacterCharacteristicForPreset struct {
	PacketId    uint32
	keyword     []byte
	base        int32
	additionnal int32
}

var simpleCharacterCharacteristicForPresetMap = make(map[uint]*simpleCharacterCharacteristicForPreset)

func GetSimpleCharacterCharacteristicForPresetNOA(instance uint) *simpleCharacterCharacteristicForPreset {
	simpleCharacterCharacteristicForPreset_, ok := simpleCharacterCharacteristicForPresetMap[instance]

	if ok {
		return simpleCharacterCharacteristicForPreset_
	}

	simpleCharacterCharacteristicForPresetMap[instance] = &simpleCharacterCharacteristicForPreset{PacketId: SimpleCharacterCharacteristicForPresetID}
	return simpleCharacterCharacteristicForPresetMap[instance]
}

func (si *simpleCharacterCharacteristicForPreset) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, si.keyword)
	utils.WriteVarInt16(buff, si.base)
	utils.WriteVarInt16(buff, si.additionnal)
}

func (si *simpleCharacterCharacteristicForPreset) Deserialize(reader *bytes.Reader) {
	si.keyword = utils.ReadUTF(reader)
	si.base = utils.ReadVarInt16(reader)
	si.additionnal = utils.ReadVarInt16(reader)
}

func (si *simpleCharacterCharacteristicForPreset) GetPacketId() uint32 {
	return si.PacketId
}

func (si *simpleCharacterCharacteristicForPreset) String() string {
	return fmt.Sprintf("packetId: %d\n", si.PacketId)
}
