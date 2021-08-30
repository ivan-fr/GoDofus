// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 14:58:10.885078 +0200 CEST m=+49.293073101

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type idol struct {
	PacketId         uint32
	id               int32
	xpBonusPercent   int32
	dropBonusPercent int32
}

var idolMap = make(map[uint]*idol)

func GetIdolNOA(instance uint) *idol {
	idol_, ok := idolMap[instance]

	if ok {
		return idol_
	}

	idolMap[instance] = &idol{PacketId: IdolID}
	return idolMap[instance]
}

func (id *idol) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt16(buff, id.id)
	utils.WriteVarInt16(buff, id.xpBonusPercent)
	utils.WriteVarInt16(buff, id.dropBonusPercent)
}

func (id *idol) Deserialize(reader *bytes.Reader) {
	id.id = utils.ReadVarInt16(reader)
	id.xpBonusPercent = utils.ReadVarInt16(reader)
	id.dropBonusPercent = utils.ReadVarInt16(reader)
}

func (id *idol) GetPacketId() uint32 {
	return id.PacketId
}

func (id *idol) String() string {
	return fmt.Sprintf("packetId: %d\n", id.PacketId)
}
