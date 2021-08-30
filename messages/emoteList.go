// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 22:31:19.5991729 +0200 CEST m=+22.980430801

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type emoteList struct {
	PacketId uint32
	emoteIds []byte
}

var emoteListMap = make(map[uint]*emoteList)

func (em *emoteList) GetNOA(instance uint) Message {
	emoteList_, ok := emoteListMap[instance]

	if ok {
		return emoteList_
	}

	emoteListMap[instance] = &emoteList{PacketId: EmoteListID}
	return emoteListMap[instance]
}

func (em *emoteList) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(em.emoteIds)))
	_ = binary.Write(buff, binary.BigEndian, em.emoteIds)
}

func (em *emoteList) Deserialize(reader *bytes.Reader) {
	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	em.emoteIds = make([]byte, len0_)
	_ = binary.Read(reader, binary.BigEndian, em.emoteIds)
}

func (em *emoteList) GetPacketId() uint32 {
	return em.PacketId
}

func (em *emoteList) String() string {
	return fmt.Sprintf("packetId: %d\n", em.PacketId)
}
