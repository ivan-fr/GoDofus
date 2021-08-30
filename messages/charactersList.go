// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 20:05:43.4423209 +0200 CEST m=+0.019882701

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type charactersList struct {
	PacketId          uint32
	bCL               *basicCharactersList
	hasStartupActions bool
}

var charactersListMap = make(map[uint]*charactersList)

func (c *charactersList) GetNOA(instance uint) Message {
	charactersList_, ok := charactersListMap[instance]

	if ok {
		return charactersList_
	}

	charactersListMap[instance] = &charactersList{PacketId: CharactersListID}
	return charactersListMap[instance]
}

func (c *charactersList) Serialize(buff *bytes.Buffer) {
	c.bCL.Serialize(buff)
	_ = binary.Write(buff, binary.BigEndian, c.hasStartupActions)
}

func (c *charactersList) Deserialize(reader *bytes.Reader) {
	c.bCL = new(basicCharactersList)
	c.bCL.Deserialize(reader)
	_ = binary.Read(reader, binary.BigEndian, &c.hasStartupActions)
}

func (c *charactersList) GetPacketId() uint32 {
	return c.PacketId
}

func (c *charactersList) String() string {
	return fmt.Sprintf("packetId: %d\n", c.PacketId)
}
