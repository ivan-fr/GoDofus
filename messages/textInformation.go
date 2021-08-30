// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 14:38:38.9109241 +0200 CEST m=+56.330106001

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type textInformation struct {
	PacketId   uint32
	msgType    byte
	msgId      int32
	parameters [][]byte
}

var textInformationMap = make(map[uint]*textInformation)

func GetTextInformationNOA(instance uint) *textInformation {
	textInformation_, ok := textInformationMap[instance]

	if ok {
		return textInformation_
	}

	textInformationMap[instance] = &textInformation{PacketId: TextInformationID}
	return textInformationMap[instance]
}

func (te *textInformation) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, te.msgType)
	utils.WriteVarInt16(buff, te.msgId)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(te.parameters)))
	for i := 0; i < len(te.parameters); i++ {
		utils.WriteUTF(buff, te.parameters[i])
	}
}

func (te *textInformation) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &te.msgType)
	te.msgId = utils.ReadVarInt16(reader)
	var len2_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len2_)
	te.parameters = nil
	for i := 0; i < int(len2_); i++ {
		te.parameters = append(te.parameters, utils.ReadUTF(reader))
	}
}

func (te *textInformation) GetPacketId() uint32 {
	return te.PacketId
}

func (te *textInformation) String() string {
	return fmt.Sprintf("packetId: %d\n", te.PacketId)
}