// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 21:47:52.0903221 +0200 CEST m=+36.540505401

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type preset struct {
	PacketId uint32
	id       uint16
}

var presetMap = make(map[uint]*preset)

func (pr *preset) GetNOA(instance uint) Message {
	preset_, ok := presetMap[instance]

	if ok {
		return preset_
	}

	presetMap[instance] = &preset{PacketId: PresetID}
	return presetMap[instance]
}

func (pr *preset) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, pr.id)
}

func (pr *preset) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &pr.id)
}

func (pr *preset) GetPacketId() uint32 {
	return pr.PacketId
}

func (pr *preset) String() string {
	return fmt.Sprintf("packetId: %d\n", pr.PacketId)
}
