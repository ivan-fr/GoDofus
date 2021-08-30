// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 14:40:16.1343363 +0200 CEST m=+28.555148601

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type spouseStatus struct {
	PacketId  uint32
	hasSpouse bool
}

var spouseStatusMap = make(map[uint]*spouseStatus)

func GetSpouseStatusNOA(instance uint) *spouseStatus {
	spouseStatus_, ok := spouseStatusMap[instance]

	if ok {
		return spouseStatus_
	}

	spouseStatusMap[instance] = &spouseStatus{PacketId: SpouseStatusID}
	return spouseStatusMap[instance]
}

func (sp *spouseStatus) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, sp.hasSpouse)
}

func (sp *spouseStatus) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &sp.hasSpouse)
}

func (sp *spouseStatus) GetPacketId() uint32 {
	return sp.PacketId
}

func (sp *spouseStatus) String() string {
	return fmt.Sprintf("packetId: %d\n", sp.PacketId)
}