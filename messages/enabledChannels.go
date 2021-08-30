// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 09:32:47.2569685 +0200 CEST m=+50.500207401

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type enabledChannels struct {
	PacketId   uint32
	channels   []byte
	disallowed []byte
}

var enabledChannelsMap = make(map[uint]*enabledChannels)

func GetEnabledChannelsNOA(instance uint) *enabledChannels {
	enabledChannels_, ok := enabledChannelsMap[instance]

	if ok {
		return enabledChannels_
	}

	enabledChannelsMap[instance] = &enabledChannels{PacketId: EnabledChannelsID}
	return enabledChannelsMap[instance]
}

func (en *enabledChannels) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(en.channels)))
	_ = binary.Write(buff, binary.BigEndian, en.channels)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(en.disallowed)))
	_ = binary.Write(buff, binary.BigEndian, en.disallowed)
}

func (en *enabledChannels) Deserialize(reader *bytes.Reader) {
	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	en.channels = make([]byte, len0_)
	_ = binary.Read(reader, binary.BigEndian, en.channels)
	var len1_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len1_)
	en.disallowed = make([]byte, len1_)
	_ = binary.Read(reader, binary.BigEndian, en.disallowed)
}

func (en *enabledChannels) GetPacketId() uint32 {
	return en.PacketId
}

func (en *enabledChannels) String() string {
	return fmt.Sprintf("packetId: %d\n", en.PacketId)
}
