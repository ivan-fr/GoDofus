// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 09:21:09.6994054 +0200 CEST m=+0.020429901

package messages

import (
	"bytes"
	"fmt"
)

type mountSet struct {
	PacketId uint32
	mC       *mountClient
}

var mountSetMap = make(map[uint]*mountSet)

func GetMountSetNOA(instance uint) *mountSet {
	mountSet_, ok := mountSetMap[instance]

	if ok {
		return mountSet_
	}

	mountSetMap[instance] = &mountSet{PacketId: MountSetID}
	return mountSetMap[instance]
}

func (m *mountSet) Serialize(buff *bytes.Buffer) {
	m.mC.Serialize(buff)
}

func (m *mountSet) Deserialize(reader *bytes.Reader) {
	m.mC = new(mountClient)
	m.mC.Deserialize(reader)
}

func (m *mountSet) GetPacketId() uint32 {
	return m.PacketId
}

func (m *mountSet) String() string {
	return fmt.Sprintf("packetId: %d\n", m.PacketId)
}