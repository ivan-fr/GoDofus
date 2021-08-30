// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 15:43:45.4168772 +0200 CEST m=+45.067695801

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type startupActionAddObject struct {
	PacketId                          uint32
	uid                               int32
	title                             []byte
	descUrl                           []byte
	pictureUrl                        []byte
	objectItemInformationWithQuantity []*objectItemInformationWithQuantity
}

var startupActionAddObjectMap = make(map[uint]*startupActionAddObject)

func GetStartupActionAddObjectNOA(instance uint) *startupActionAddObject {
	startupActionAddObject_, ok := startupActionAddObjectMap[instance]

	if ok {
		return startupActionAddObject_
	}

	startupActionAddObjectMap[instance] = &startupActionAddObject{PacketId: StartupActionAddObjectID}
	return startupActionAddObjectMap[instance]
}

func (st *startupActionAddObject) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, st.uid)
	utils.WriteUTF(buff, st.title)
	utils.WriteUTF(buff, st.descUrl)
	utils.WriteUTF(buff, st.pictureUrl)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(st.objectItemInformationWithQuantity)))
	for i := 0; i < len(st.objectItemInformationWithQuantity); i++ {
		st.objectItemInformationWithQuantity[i].Serialize(buff)
	}
}

func (st *startupActionAddObject) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &st.uid)
	st.title = utils.ReadUTF(reader)
	st.descUrl = utils.ReadUTF(reader)
	st.pictureUrl = utils.ReadUTF(reader)
	var len4_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len4_)
	st.objectItemInformationWithQuantity = nil
	for i := 0; i < int(len4_); i++ {
		aMessage4 := new(objectItemInformationWithQuantity)
		aMessage4.Deserialize(reader)
		st.objectItemInformationWithQuantity = append(st.objectItemInformationWithQuantity, aMessage4)
	}
}

func (st *startupActionAddObject) GetPacketId() uint32 {
	return st.PacketId
}

func (st *startupActionAddObject) String() string {
	return fmt.Sprintf("packetId: %d\n", st.PacketId)
}
