// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 20:50:35.0532152 +0200 CEST m=+0.020007601

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type subEntity struct {
	PacketId             uint32
	bindingPointCategory byte
	bindingPointIndex    byte
	subEntityLook        *entityLook
}

var subEntityMap = make(map[uint]*subEntity)

func GetSubEntityNOA(instance uint) *subEntity {
	subEntity_, ok := subEntityMap[instance]

	if ok {
		return subEntity_
	}

	subEntityMap[instance] = &subEntity{PacketId: SubEntityID}
	return subEntityMap[instance]
}

func (s *subEntity) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, s.bindingPointCategory)
	_ = binary.Write(buff, binary.BigEndian, s.bindingPointIndex)
	s.subEntityLook.Serialize(buff)
}

func (s *subEntity) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &s.bindingPointCategory)
	_ = binary.Read(reader, binary.BigEndian, &s.bindingPointIndex)
	s.subEntityLook = new(entityLook)
	s.subEntityLook.Deserialize(reader)
}

func (s *subEntity) GetPacketId() uint32 {
	return s.PacketId
}

func (s *subEntity) String() string {
	return fmt.Sprintf("packetId: %d\n", s.PacketId)
}