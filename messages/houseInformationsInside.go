// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 15:15:44.2038952 +0200 CEST m=+103.764523801

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type houseInformationsInside struct {
	PacketId          uint32
	houseInformations *houseInformations
	item              *item
	worldX            uint16
	worldY            uint16
}

var houseInformationsInsideMap = make(map[uint]*houseInformationsInside)

func GetHouseInformationsInsideNOA(instance uint) *houseInformationsInside {
	houseInformationsInside_, ok := houseInformationsInsideMap[instance]

	if ok {
		return houseInformationsInside_
	}

	houseInformationsInsideMap[instance] = &houseInformationsInside{PacketId: HouseInformationsInsideID}
	return houseInformationsInsideMap[instance]
}

func (ho *houseInformationsInside) Serialize(buff *bytes.Buffer) {
	ho.houseInformations.Serialize(buff)
	ho.item.Serialize(buff)
	_ = binary.Write(buff, binary.BigEndian, ho.worldX)
	_ = binary.Write(buff, binary.BigEndian, ho.worldY)
}

func (ho *houseInformationsInside) Deserialize(reader *bytes.Reader) {
	ho.houseInformations = new(houseInformations)
	ho.houseInformations.Deserialize(reader)
	ho.item = new(item)
	ho.item.Deserialize(reader)
	_ = binary.Read(reader, binary.BigEndian, &ho.worldX)
	_ = binary.Read(reader, binary.BigEndian, &ho.worldY)
}

func (ho *houseInformationsInside) GetPacketId() uint32 {
	return ho.PacketId
}

func (ho *houseInformationsInside) String() string {
	return fmt.Sprintf("packetId: %d\n", ho.PacketId)
}