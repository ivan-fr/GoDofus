// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 15:35:17.0969196 +0200 CEST m=+58.663678001

package messages

import (
	"bytes"
	"fmt"
)

type houseGuildedInformations struct {
	PacketId                  uint32
	houseInstanceInformations *houseInstanceInformations
	guildInformations         *guildInformations
}

var houseGuildedInformationsMap = make(map[uint]*houseGuildedInformations)

func GetHouseGuildedInformationsNOA(instance uint) *houseGuildedInformations {
	houseGuildedInformations_, ok := houseGuildedInformationsMap[instance]

	if ok {
		return houseGuildedInformations_
	}

	houseGuildedInformationsMap[instance] = &houseGuildedInformations{PacketId: HouseGuildedInformationsID}
	return houseGuildedInformationsMap[instance]
}

func (ho *houseGuildedInformations) Serialize(buff *bytes.Buffer) {
	ho.houseInstanceInformations.Serialize(buff)
	ho.guildInformations.Serialize(buff)
}

func (ho *houseGuildedInformations) Deserialize(reader *bytes.Reader) {
	ho.houseInstanceInformations = new(houseInstanceInformations)
	ho.houseInstanceInformations.Deserialize(reader)
	ho.guildInformations = new(guildInformations)
	ho.guildInformations.Deserialize(reader)
}

func (ho *houseGuildedInformations) GetPacketId() uint32 {
	return ho.PacketId
}

func (ho *houseGuildedInformations) String() string {
	return fmt.Sprintf("packetId: %d\n", ho.PacketId)
}