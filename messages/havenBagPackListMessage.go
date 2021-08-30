// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 16:04:38.2953863 +0200 CEST m=+15.000551001

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type havenBagPackListMessage struct {
	PacketId uint32
	packIds  []byte
}

var havenBagPackListMessageMap = make(map[uint]*havenBagPackListMessage)

func (ha *havenBagPackListMessage) GetNOA(instance uint) Message {
	havenBagPackListMessage_, ok := havenBagPackListMessageMap[instance]

	if ok {
		return havenBagPackListMessage_
	}

	havenBagPackListMessageMap[instance] = &havenBagPackListMessage{PacketId: HavenBagPackListMessageID}
	return havenBagPackListMessageMap[instance]
}

func (ha *havenBagPackListMessage) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(ha.packIds)))
	_ = binary.Write(buff, binary.BigEndian, ha.packIds)
}

func (ha *havenBagPackListMessage) Deserialize(reader *bytes.Reader) {
	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	ha.packIds = make([]byte, len0_)
	_ = binary.Read(reader, binary.BigEndian, ha.packIds)
}

func (ha *havenBagPackListMessage) GetPacketId() uint32 {
	return ha.PacketId
}

func (ha *havenBagPackListMessage) String() string {
	return fmt.Sprintf("packetId: %d\n", ha.PacketId)
}
