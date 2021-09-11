// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-09-11 12:30:10.4989814 +0200 CEST m=+24.507856501

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type MapObstacle struct {
	PacketId       uint32
	obstacleCellId int32
	state          byte
}

var MapObstacleMap = make(map[uint]*MapObstacle)

func (Ma *MapObstacle) GetNOA(instance uint) Message {
	MapObstacle_, ok := MapObstacleMap[instance]

	if ok {
		return MapObstacle_
	}

	MapObstacleMap[instance] = &MapObstacle{PacketId: MapObstacleID}
	return MapObstacleMap[instance]
}

func (Ma *MapObstacle) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt16(buff, Ma.obstacleCellId)
	_ = binary.Write(buff, binary.BigEndian, Ma.state)
}

func (Ma *MapObstacle) Deserialize(reader *bytes.Reader) {
	Ma.obstacleCellId = utils.ReadVarInt16(reader)
	_ = binary.Read(reader, binary.BigEndian, &Ma.state)
}

func (Ma *MapObstacle) GetPacketId() uint32 {
	return Ma.PacketId
}

func (Ma *MapObstacle) String() string {
	return fmt.Sprintf("packetId: %d\n", Ma.PacketId)
}