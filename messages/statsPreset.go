// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 22:11:40.8075143 +0200 CEST m=+51.782592501

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type statsPreset struct {
	PacketId uint32
	preset   *preset
	sCCFPs   []*simpleCharacterCharacteristicForPreset
}

var statsPresetMap = make(map[uint]*statsPreset)

func (st *statsPreset) GetNOA(instance uint) Message {
	statsPreset_, ok := statsPresetMap[instance]

	if ok {
		return statsPreset_
	}

	statsPresetMap[instance] = &statsPreset{PacketId: StatsPresetID}
	return statsPresetMap[instance]
}

func (st *statsPreset) Serialize(buff *bytes.Buffer) {
	st.preset.Serialize(buff)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(st.sCCFPs)))
	for i := 0; i < len(st.sCCFPs); i++ {
		st.sCCFPs[i].Serialize(buff)
	}
}

func (st *statsPreset) Deserialize(reader *bytes.Reader) {
	st.preset = new(preset)
	st.preset.Deserialize(reader)
	var len1_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len1_)
	st.sCCFPs = nil
	for i := 0; i < int(len1_); i++ {
		aMessage1 := new(simpleCharacterCharacteristicForPreset)
		aMessage1.Deserialize(reader)
		st.sCCFPs = append(st.sCCFPs, aMessage1)
	}
}

func (st *statsPreset) GetPacketId() uint32 {
	return st.PacketId
}

func (st *statsPreset) String() string {
	return fmt.Sprintf("packetId: %d\n", st.PacketId)
}
