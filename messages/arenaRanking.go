// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 15:46:52.5365556 +0200 CEST m=+22.109944001

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type arenaRanking struct {
	PacketId uint32
	rank     int32
	bestRank int32
}

var arenaRankingMap = make(map[uint]*arenaRanking)

func GetArenaRankingNOA(instance uint) *arenaRanking {
	arenaRanking_, ok := arenaRankingMap[instance]

	if ok {
		return arenaRanking_
	}

	arenaRankingMap[instance] = &arenaRanking{PacketId: ArenaRankingID}
	return arenaRankingMap[instance]
}

func (ar *arenaRanking) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt16(buff, ar.rank)
	utils.WriteVarInt16(buff, ar.bestRank)
}

func (ar *arenaRanking) Deserialize(reader *bytes.Reader) {
	ar.rank = utils.ReadVarInt16(reader)
	ar.bestRank = utils.ReadVarInt16(reader)
}

func (ar *arenaRanking) GetPacketId() uint32 {
	return ar.PacketId
}

func (ar *arenaRanking) String() string {
	return fmt.Sprintf("packetId: %d\n", ar.PacketId)
}