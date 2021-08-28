// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 17:16:46.256238 +0200 CEST m=+0.019798401

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type serverOptionalFeatures struct {
	PacketId uint32
	features []byte
}

var serverOptionalFeaturesMap = make(map[uint]*serverOptionalFeatures)

func GetServerOptionalFeaturesNOA(instance uint) *serverOptionalFeatures {
	serverOptionalFeatures_, ok := serverOptionalFeaturesMap[instance]

	if ok {
		return serverOptionalFeatures_
	}

	serverOptionalFeaturesMap[instance] = &serverOptionalFeatures{PacketId: ServerOptionalFeaturesID}
	return serverOptionalFeaturesMap[instance]
}

func (s *serverOptionalFeatures) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(s.features)))
	_ = binary.Write(buff, binary.BigEndian, s.features)
}

func (s *serverOptionalFeatures) Deserialize(reader *bytes.Reader) {
	var len_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len_)

	s.features = make([]byte, len_)
	_ = binary.Read(reader, binary.BigEndian, s.features)
}

func (s *serverOptionalFeatures) GetPacketId() uint32 {
	return s.PacketId
}

func (s *serverOptionalFeatures) String() string {
	return fmt.Sprintf("packetId: %d\n", s.PacketId)
}
