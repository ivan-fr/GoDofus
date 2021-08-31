// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 23:15:34.137258 +0200 CEST m=+14.060220101

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type gameContextCreate struct {
	PacketId uint32
	context  byte
}

var gameContextCreateMap = make(map[uint]*gameContextCreate)

func (ga *gameContextCreate) GetNOA(instance uint) Message {
	gameContextCreate_, ok := gameContextCreateMap[instance]

	if ok {
		return gameContextCreate_
	}

	gameContextCreateMap[instance] = &gameContextCreate{PacketId: GameContextCreateID}
	return gameContextCreateMap[instance]
}

func (ga *gameContextCreate) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, ga.context)
}

func (ga *gameContextCreate) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &ga.context)
}

func (ga *gameContextCreate) GetPacketId() uint32 {
	return ga.PacketId
}

func (ga *gameContextCreate) String() string {
	return fmt.Sprintf("packetId: %d\n", ga.PacketId)
}
