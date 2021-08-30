// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 22:21:16.4959975 +0200 CEST m=+64.475734001

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type fullStatsPreset struct {
	PacketId                         uint32
	preset                           *preset
	characterCharacteristicForPreset []*characterCharacteristicForPreset
}

var fullStatsPresetMap = make(map[uint]*fullStatsPreset)

func (fu *fullStatsPreset) GetNOA(instance uint) Message {
	fullStatsPreset_, ok := fullStatsPresetMap[instance]

	if ok {
		return fullStatsPreset_
	}

	fullStatsPresetMap[instance] = &fullStatsPreset{PacketId: FullStatsPresetID}
	return fullStatsPresetMap[instance]
}

func (fu *fullStatsPreset) Serialize(buff *bytes.Buffer) {
	fu.preset.Serialize(buff)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(fu.characterCharacteristicForPreset)))
	for i := 0; i < len(fu.characterCharacteristicForPreset); i++ {
		fu.characterCharacteristicForPreset[i].Serialize(buff)
	}
}

func (fu *fullStatsPreset) Deserialize(reader *bytes.Reader) {
	fu.preset = new(preset)
	fu.preset.Deserialize(reader)
	var len1_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len1_)
	fu.characterCharacteristicForPreset = nil
	for i := 0; i < int(len1_); i++ {
		aMessage1 := new(characterCharacteristicForPreset)
		aMessage1.Deserialize(reader)
		fu.characterCharacteristicForPreset = append(fu.characterCharacteristicForPreset, aMessage1)
	}
}

func (fu *fullStatsPreset) GetPacketId() uint32 {
	return fu.PacketId
}

func (fu *fullStatsPreset) String() string {
	return fmt.Sprintf("packetId: %d\n", fu.PacketId)
}
