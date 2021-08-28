// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-25 10:47:08.6457558 +0200 CEST m=+0.002560401

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type gameServerInformations struct {
	PacketId        uint32
	isMonoAccount   bool
	isSelectable    bool
	id              int32
	type_           byte
	status          byte
	completion      byte
	charactersCount byte
	charactersSlots byte
	date            float64
}

var gameServerInformationsMap = make(map[uint]*gameServerInformations)

func GetGameServerInformationsNOA(instance uint) *gameServerInformations {
	gameServerInformations_, ok := gameServerInformationsMap[instance]

	if ok {
		return gameServerInformations_
	}

	gameServerInformationsMap[instance] = &gameServerInformations{PacketId: GameServerInformationID}
	return gameServerInformationsMap[instance]
}

func (g *gameServerInformations) Serialize(buff *bytes.Buffer) {
	var box uint32
	box = utils.SetFlag(box, 0, g.isMonoAccount)
	box = utils.SetFlag(box, 1, g.isSelectable)

	_ = binary.Write(buff, binary.BigEndian, byte(box))

	utils.WriteVarShort(buff, g.id)
	_ = binary.Write(buff, binary.BigEndian, g.type_)
	_ = binary.Write(buff, binary.BigEndian, g.status)
	_ = binary.Write(buff, binary.BigEndian, g.completion)
	_ = binary.Write(buff, binary.BigEndian, g.charactersCount)
	_ = binary.Write(buff, binary.BigEndian, g.charactersSlots)
	_ = binary.Write(buff, binary.BigEndian, g.date)
}

func (g *gameServerInformations) Deserialize(reader *bytes.Reader) {
	var box byte
	_ = binary.Read(reader, binary.BigEndian, &box)

	g.isMonoAccount = utils.GetFlag(uint32(box), 0)
	g.isSelectable = utils.GetFlag(uint32(box), 1)

	g.id = utils.ReadVarInt16(reader)
	_ = binary.Read(reader, binary.BigEndian, &g.type_)
	_ = binary.Read(reader, binary.BigEndian, &g.status)
	_ = binary.Read(reader, binary.BigEndian, &g.completion)
	_ = binary.Read(reader, binary.BigEndian, &g.charactersCount)
	_ = binary.Read(reader, binary.BigEndian, &g.charactersSlots)
	_ = binary.Read(reader, binary.BigEndian, &g.date)
}

func (g *gameServerInformations) GetPacketId() uint32 {
	return g.PacketId
}

func (g *gameServerInformations) String() string {
	return fmt.Sprintf("PacketId: %d\n", g.PacketId)
}
