// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 15:48:29.2491238 +0200 CEST m=+57.220739501

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type arenaLeagueRanking struct {
	PacketId          uint32
	rank              int32
	leagueId          int32
	leaguePoints      int32
	totalLeaguePoints int32
	ladderPosition    int32
}

var arenaLeagueRankingMap = make(map[uint]*arenaLeagueRanking)

func (ar *arenaLeagueRanking) GetNOA(instance uint) Message {
	arenaLeagueRanking_, ok := arenaLeagueRankingMap[instance]

	if ok {
		return arenaLeagueRanking_
	}

	arenaLeagueRankingMap[instance] = &arenaLeagueRanking{PacketId: ArenaLeagueRankingID}
	return arenaLeagueRankingMap[instance]
}

func (ar *arenaLeagueRanking) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt16(buff, ar.rank)
	utils.WriteVarInt16(buff, ar.leagueId)
	utils.WriteVarInt16(buff, ar.leaguePoints)
	utils.WriteVarInt16(buff, ar.totalLeaguePoints)
	_ = binary.Write(buff, binary.BigEndian, ar.ladderPosition)
}

func (ar *arenaLeagueRanking) Deserialize(reader *bytes.Reader) {
	ar.rank = utils.ReadVarInt16(reader)
	ar.leagueId = utils.ReadVarInt16(reader)
	ar.leaguePoints = utils.ReadVarInt16(reader)
	ar.totalLeaguePoints = utils.ReadVarInt16(reader)
	_ = binary.Read(reader, binary.BigEndian, &ar.ladderPosition)
}

func (ar *arenaLeagueRanking) GetPacketId() uint32 {
	return ar.PacketId
}

func (ar *arenaLeagueRanking) String() string {
	return fmt.Sprintf("packetId: %d\n", ar.PacketId)
}
