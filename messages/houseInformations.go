// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 15:07:14.9313142 +0200 CEST m=+38.948813801

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type houseInformations struct {
	PacketId uint32
	houseId  int32
	modelId  int32
}

var houseInformationsMap = make(map[uint]*houseInformations)

func GetHouseInformationsNOA(instance uint) *houseInformations {
	houseInformations_, ok := houseInformationsMap[instance]

	if ok {
		return houseInformations_
	}

	houseInformationsMap[instance] = &houseInformations{PacketId: HouseInformationsID}
	return houseInformationsMap[instance]
}

func (ho *houseInformations) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt32(buff, ho.houseId)
	utils.WriteVarInt16(buff, ho.modelId)
}

func (ho *houseInformations) Deserialize(reader *bytes.Reader) {
	ho.houseId = utils.ReadVarInt32(reader)
	ho.modelId = utils.ReadVarInt16(reader)
}

func (ho *houseInformations) GetPacketId() uint32 {
	return ho.PacketId
}

func (ho *houseInformations) String() string {
	return fmt.Sprintf("packetId: %d\n", ho.PacketId)
}