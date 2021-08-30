// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 22:49:51.2847445 +0200 CEST m=+20.730216601

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type skillActionDescriptionTimed struct {
	PacketId               uint32
	skillActionDescription *skillActionDescription
	time                   byte
}

var skillActionDescriptionTimedMap = make(map[uint]*skillActionDescriptionTimed)

func (sk *skillActionDescriptionTimed) GetNOA(instance uint) Message {
	skillActionDescriptionTimed_, ok := skillActionDescriptionTimedMap[instance]

	if ok {
		return skillActionDescriptionTimed_
	}

	skillActionDescriptionTimedMap[instance] = &skillActionDescriptionTimed{PacketId: SkillActionDescriptionTimedID}
	return skillActionDescriptionTimedMap[instance]
}

func (sk *skillActionDescriptionTimed) Serialize(buff *bytes.Buffer) {
	sk.skillActionDescription.Serialize(buff)
	_ = binary.Write(buff, binary.BigEndian, sk.time)
}

func (sk *skillActionDescriptionTimed) Deserialize(reader *bytes.Reader) {
	sk.skillActionDescription = new(skillActionDescription)
	sk.skillActionDescription.Deserialize(reader)
	_ = binary.Read(reader, binary.BigEndian, &sk.time)
}

func (sk *skillActionDescriptionTimed) GetPacketId() uint32 {
	return sk.PacketId
}

func (sk *skillActionDescriptionTimed) String() string {
	return fmt.Sprintf("packetId: %d\n", sk.PacketId)
}
