// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 08:53:52.4042438 +0200 CEST m=+40.533736501

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type prismSubareaEmptyInfo struct {
	PacketId   uint32
	subAreaId  int32
	allianceId int32
}

var prismSubareaEmptyInfoMap = make(map[uint]*prismSubareaEmptyInfo)

func (pr *prismSubareaEmptyInfo) GetNOA(instance uint) Message {
	prismSubareaEmptyInfo_, ok := prismSubareaEmptyInfoMap[instance]

	if ok {
		return prismSubareaEmptyInfo_
	}

	prismSubareaEmptyInfoMap[instance] = &prismSubareaEmptyInfo{PacketId: PrismSubareaEmptyInfoID}
	return prismSubareaEmptyInfoMap[instance]
}

func (pr *prismSubareaEmptyInfo) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt16(buff, pr.subAreaId)
	utils.WriteVarInt32(buff, pr.allianceId)
}

func (pr *prismSubareaEmptyInfo) Deserialize(reader *bytes.Reader) {
	pr.subAreaId = utils.ReadVarInt16(reader)
	pr.allianceId = utils.ReadVarInt32(reader)
}

func (pr *prismSubareaEmptyInfo) GetPacketId() uint32 {
	return pr.PacketId
}

func (pr *prismSubareaEmptyInfo) String() string {
	return fmt.Sprintf("packetId: %d\n", pr.PacketId)
}
