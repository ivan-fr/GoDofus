// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 17:20:19.3370255 +0200 CEST m=+0.019788001

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type haapiSession struct {
	PacketId uint32
	key      []byte
	type_    byte
}

var haapiSessionMap = make(map[uint]*haapiSession)

func (h *haapiSession) GetNOA(instance uint) Message {
	haapiSession_, ok := haapiSessionMap[instance]

	if ok {
		return haapiSession_
	}

	haapiSessionMap[instance] = &haapiSession{PacketId: HaapiSessionID}
	return haapiSessionMap[instance]
}

func (h *haapiSession) Serialize(buff *bytes.Buffer) {
	utils.WriteUTF(buff, h.key)
	_ = binary.Write(buff, binary.BigEndian, h.type_)
}

func (h *haapiSession) Deserialize(reader *bytes.Reader) {
	h.key = utils.ReadUTF(reader)
	_ = binary.Read(reader, binary.BigEndian, &h.type_)
}

func (h *haapiSession) GetPacketId() uint32 {
	return h.PacketId
}

func (h *haapiSession) String() string {
	return fmt.Sprintf("packetId: %d\n", h.PacketId)
}
