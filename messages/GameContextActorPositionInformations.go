// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-09-11 12:03:54.7322858 +0200 CEST m=+52.858237801

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type GameContextActorPositionInformations struct {
	PacketId                      uint32
	contextualId                  float64
	EntityDispositionInformations *item
}

var GameContextActorPositionInformationsMap = make(map[uint]*GameContextActorPositionInformations)

func (Ga *GameContextActorPositionInformations) GetNOA(instance uint) Message {
	GameContextActorPositionInformations_, ok := GameContextActorPositionInformationsMap[instance]

	if ok {
		return GameContextActorPositionInformations_
	}

	GameContextActorPositionInformationsMap[instance] = &GameContextActorPositionInformations{PacketId: GameContextActorPositionInformationsID}
	return GameContextActorPositionInformationsMap[instance]
}

func (Ga *GameContextActorPositionInformations) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, Ga.contextualId)
	Ga.EntityDispositionInformations.Serialize(buff)
}

func (Ga *GameContextActorPositionInformations) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &Ga.contextualId)
	Ga.EntityDispositionInformations = new(item)
	Ga.EntityDispositionInformations.Deserialize(reader)
}

func (Ga *GameContextActorPositionInformations) GetPacketId() uint32 {
	return Ga.PacketId
}

func (Ga *GameContextActorPositionInformations) String() string {
	return fmt.Sprintf("packetId: %d\n", Ga.PacketId)
}