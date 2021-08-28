// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 20:43:11.2229317 +0200 CEST m=+0.020264401

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type entityLook struct {
	PacketId      uint32
	bonesId       int32
	skins         []int32
	indexedColors []int32
	scales        []int32
	subentities   []*subEntity
}

var entityLookMap = make(map[uint]*entityLook)

func GetEntityLookNOA(instance uint) *entityLook {
	entityLook_, ok := entityLookMap[instance]

	if ok {
		return entityLook_
	}

	entityLookMap[instance] = &entityLook{PacketId: EntityLookID}
	return entityLookMap[instance]
}

func (e *entityLook) Serialize(buff *bytes.Buffer) {
	utils.WriteVarShort(buff, e.bonesId)

	_ = binary.Write(buff, binary.BigEndian, uint16(len(e.skins)))
	for i := 0; i < len(e.skins); i++ {
		utils.WriteVarShort(buff, e.skins[i])
	}

	_ = binary.Write(buff, binary.BigEndian, uint16(len(e.indexedColors)))
	_ = binary.Write(buff, binary.BigEndian, e.indexedColors)

	_ = binary.Write(buff, binary.BigEndian, uint16(len(e.scales)))
	for i := 0; i < len(e.skins); i++ {
		utils.WriteVarShort(buff, e.scales[i])
	}

	_ = binary.Write(buff, binary.BigEndian, uint16(len(e.subentities)))
	for i := 0; i < len(e.skins); i++ {
		e.subentities[i].Serialize(buff)
	}
}

func (e *entityLook) Deserialize(reader *bytes.Reader) {
	e.bonesId = utils.ReadVarInt16(reader)

	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	for i := 0; i < int(len0_); i++ {
		e.skins[i] = utils.ReadVarInt16(reader)
	}

	var len1_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len1_)
	e.indexedColors = make([]int32, len1_)
	_ = binary.Read(reader, binary.BigEndian, e.indexedColors)

	var len2_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len2_)
	for i := 0; i < int(len2_); i++ {
		e.scales[i] = utils.ReadVarInt32(reader)
	}

	var len3_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len3_)
	for i := 0; i < int(len3_); i++ {
		aSubEntity := new(subEntity)
		aSubEntity.Deserialize(reader)
		e.subentities = append(e.subentities, aSubEntity)
	}
}

func (e *entityLook) GetPacketId() uint32 {
	return e.PacketId
}

func (e *entityLook) String() string {
	return fmt.Sprintf("packetId: %d\n", e.PacketId)
}
