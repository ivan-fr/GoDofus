// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 08:50:52.5502926 +0200 CEST m=+0.019885401

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type emoteAdd struct {
	PacketId uint32
	emoteId  byte
}

var emoteAddMap = make(map[uint]*emoteAdd)

func GetEmoteAddNOA(instance uint) *emoteAdd {
	emoteAdd_, ok := emoteAddMap[instance]

	if ok {
		return emoteAdd_
	}

	emoteAddMap[instance] = &emoteAdd{PacketId: EmoteAddID}
	return emoteAddMap[instance]
}

func (e *emoteAdd) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, e.emoteId)
}

func (e *emoteAdd) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &e.emoteId)
}

func (e *emoteAdd) GetPacketId() uint32 {
	return e.PacketId
}

func (e *emoteAdd) String() string {
	return fmt.Sprintf("packetId: %d\n", e.PacketId)
}
