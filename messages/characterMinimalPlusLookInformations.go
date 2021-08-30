// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 20:33:01.5921461 +0200 CEST m=+0.020621501

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type characterMinimalPlusLookInformations struct {
	PacketId uint32
	cMI      *characterMinimalInformations
	eL       *entityLook
	breed    byte
}

var characterMinimalPlusLookInformationsMap = make(map[uint]*characterMinimalPlusLookInformations)

func (c *characterMinimalPlusLookInformations) GetNOA(instance uint) Message {
	characterMinimalPlusLookInformations_, ok := characterMinimalPlusLookInformationsMap[instance]

	if ok {
		return characterMinimalPlusLookInformations_
	}

	characterMinimalPlusLookInformationsMap[instance] = &characterMinimalPlusLookInformations{PacketId: CharacterMinimalPlusLookInformationsID}
	return characterMinimalPlusLookInformationsMap[instance]
}

func (c *characterMinimalPlusLookInformations) Serialize(buff *bytes.Buffer) {
	c.cMI.Serialize(buff)
	c.eL.Serialize(buff)
	_ = binary.Write(buff, binary.BigEndian, c.breed)
}

func (c *characterMinimalPlusLookInformations) Deserialize(reader *bytes.Reader) {
	c.cMI = new(characterMinimalInformations)
	c.cMI.Deserialize(reader)
	c.eL = new(entityLook)
	c.eL.Deserialize(reader)
	_ = binary.Read(reader, binary.BigEndian, &c.breed)
}

func (c *characterMinimalPlusLookInformations) GetPacketId() uint32 {
	return c.PacketId
}

func (c *characterMinimalPlusLookInformations) String() string {
	return fmt.Sprintf("packetId: %d\n", c.PacketId)
}
