// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-25 13:37:29.6275734 +0200 CEST m=+0.003123301

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type rawData struct {
	PacketId uint32
	content  []byte
}

var rawData_ = &rawData{PacketId: RawDataID}

func GetRawDataNOA() *rawData {
	return rawData_
}

func (r *rawData) Serialize(buff *bytes.Buffer) {

}

func (r *rawData) Deserialize(reader *bytes.Reader) {
	contentLen := utils.ReadVarInt32(reader)
	r.content = make([]byte, contentLen)
	_ = binary.Read(reader, binary.BigEndian, r.content)
}

func (r *rawData) GetPacketId() uint32 {
	return r.PacketId
}

func (r *rawData) String() string {
	return fmt.Sprintf("packetId: %d\n", r.PacketId)
}