// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 22:27:45.7176231 +0200 CEST m=+32.141427201

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type channelEnabling struct {
	PacketId uint32
	channel  byte
	enable   bool
}

var channelEnablingMap = make(map[uint]*channelEnabling)

func (ch *channelEnabling) GetNOA(instance uint) Message {
	channelEnabling_, ok := channelEnablingMap[instance]

	if ok {
		return channelEnabling_
	}

	channelEnablingMap[instance] = &channelEnabling{PacketId: ChannelEnablingID}
	return channelEnablingMap[instance]
}

func (ch *channelEnabling) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, ch.channel)
	_ = binary.Write(buff, binary.BigEndian, ch.enable)
}

func (ch *channelEnabling) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &ch.channel)
	_ = binary.Read(reader, binary.BigEndian, &ch.enable)
}

func (ch *channelEnabling) GetPacketId() uint32 {
	return ch.PacketId
}

func (ch *channelEnabling) String() string {
	return fmt.Sprintf("packetId: %d\n", ch.PacketId)
}
