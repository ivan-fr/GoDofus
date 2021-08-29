// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 11:50:30.1265896 +0200 CEST m=+0.020364101

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type objectEffectCreature struct {
	PacketId        uint32
	oE              *objectEffect
	monsterFamilyId int32
}

var objectEffectCreatureMap = make(map[uint]*objectEffectCreature)

func GetObjectEffectCreatureNOA(instance uint) *objectEffectCreature {
	objectEffectCreature_, ok := objectEffectCreatureMap[instance]

	if ok {
		return objectEffectCreature_
	}

	objectEffectCreatureMap[instance] = &objectEffectCreature{PacketId: ObjectEffectCreatureID}
	return objectEffectCreatureMap[instance]
}

func (o *objectEffectCreature) Serialize(buff *bytes.Buffer) {
	o.oE.Serialize(buff)
	utils.WriteVarInt16(buff, o.monsterFamilyId)
}

func (o *objectEffectCreature) Deserialize(reader *bytes.Reader) {
	o.oE = new(objectEffect)
	o.oE.Deserialize(reader)
	o.monsterFamilyId = utils.ReadVarInt16(reader)
}

func (o *objectEffectCreature) GetPacketId() uint32 {
	return o.PacketId
}

func (o *objectEffectCreature) String() string {
	return fmt.Sprintf("packetId: %d\n", o.PacketId)
}
