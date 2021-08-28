// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 19:55:17.6573279 +0200 CEST m=+0.019756701

package messages

import (
	"bytes"
	"fmt"
)

type haapiApiKeyRequest struct {
	PacketId uint32
}

var haapiApiKeyRequestMap = make(map[uint]*haapiApiKeyRequest)

func GetHaapiApiKeyRequestNOA(instance uint) *haapiApiKeyRequest {
	haapiApiKeyRequest_, ok := haapiApiKeyRequestMap[instance]

	if ok {
		return haapiApiKeyRequest_
	}

	haapiApiKeyRequestMap[instance] = &haapiApiKeyRequest{PacketId: HaapiApiKeyRequestID}
	return haapiApiKeyRequestMap[instance]
}

func (h *haapiApiKeyRequest) Serialize(buff *bytes.Buffer) {

}

func (h *haapiApiKeyRequest) Deserialize(reader *bytes.Reader) {

}

func (h *haapiApiKeyRequest) GetPacketId() uint32 {
	return h.PacketId
}

func (h *haapiApiKeyRequest) String() string {
	return fmt.Sprintf("packetId: %d\n", h.PacketId)
}
