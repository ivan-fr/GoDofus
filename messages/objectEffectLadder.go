// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 11:53:06.1182268 +0200 CEST m=+0.019984801

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type objectEffectLadder struct {
	PacketId     uint32
	oEC          *objectEffectCreature
	monsterCount int32
}

var objectEffectLadderMap = make(map[uint]*objectEffectLadder)

func GetObjectEffectLadderNOA(instance uint) *objectEffectLadder {
	objectEffectLadder_, ok := objectEffectLadderMap[instance]

	if ok {
		return objectEffectLadder_
	}

	objectEffectLadderMap[instance] = &objectEffectLadder{PacketId: ObjectEffectLadderID}
	return objectEffectLadderMap[instance]
}

func (o *objectEffectLadder) Serialize(buff *bytes.Buffer) {
	o.oEC.Serialize(buff)
	utils.WriteVarInt32(buff, o.monsterCount)
}

func (o *objectEffectLadder) Deserialize(reader *bytes.Reader) {
	o.oEC = new(objectEffectCreature)
	o.oEC.Deserialize(reader)
	o.monsterCount = utils.ReadVarInt32(reader)
}

func (o *objectEffectLadder) GetPacketId() uint32 {
	return o.PacketId
}

func (o *objectEffectLadder) String() string {
	return fmt.Sprintf("packetId: %d\n", o.PacketId)
}
