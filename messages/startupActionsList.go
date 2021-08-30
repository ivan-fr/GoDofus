// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 22:47:42.3409179 +0200 CEST m=+36.238866701

package messages

import (
	"bytes"
	"fmt"
)

type startupActionsList struct {
	PacketId                uint32
	startupActionAddObject0 []*startupActionAddObject
}

var startupActionsListMap = make(map[uint]*startupActionsList)

func (st *startupActionsList) GetNOA(instance uint) Message {
	startupActionsList_, ok := startupActionsListMap[instance]

	if ok {
		return startupActionsList_
	}

	startupActionsListMap[instance] = &startupActionsList{PacketId: StartupActionsListID}
	return startupActionsListMap[instance]
}

func (st *startupActionsList) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(st.startupActionAddObject0)))
	for i := 0; i < len(st.startupActionAddObject0); i++ {
		st.startupActionAddObject0[i].Serialize(buff)
	}
}

func (st *startupActionsList) Deserialize(reader *bytes.Reader) {
	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	st.startupActionAddObject0 = nil
	for i := 0; i < int(len0_); i++ {
		aMessage0 := new(startupActionAddObject)
		aMessage0.Deserialize(reader)
		st.startupActionAddObject0 = append(st.startupActionAddObject0, aMessage0)
	}
}

func (st *startupActionsList) GetPacketId() uint32 {
	return st.PacketId
}

func (st *startupActionsList) String() string {
	return fmt.Sprintf("packetId: %d\n", st.PacketId)
}
