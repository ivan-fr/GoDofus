// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-25 13:28:37.6022356 +0200 CEST m=+0.002588901

package messages

import (
	"bytes"
	"fmt"
)

type helloGame struct {
	PacketId uint32
}

var helloGame_ = &helloGame{PacketId: HelloGameID}

func GetHelloGameNOA() *helloGame {
	return helloGame_
}

func (h *helloGame) Serialize(buff *bytes.Buffer) {

}

func (h *helloGame) Deserialize(reader *bytes.Reader) {

}

func (h *helloGame) GetPacketId() uint32 {
	return h.PacketId
}

func (h *helloGame) String() string {
	return fmt.Sprintf("packetId: %d\n", h.PacketId)
}
