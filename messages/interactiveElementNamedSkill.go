// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 22:46:44.2854082 +0200 CEST m=+22.048958801

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type interactiveElementNamedSkill struct {
	PacketId                uint32
	interactiveElementSkill *interactiveElementSkill
	nameId                  int32
}

var interactiveElementNamedSkillMap = make(map[uint]*interactiveElementNamedSkill)

func (in *interactiveElementNamedSkill) GetNOA(instance uint) Message {
	interactiveElementNamedSkill_, ok := interactiveElementNamedSkillMap[instance]

	if ok {
		return interactiveElementNamedSkill_
	}

	interactiveElementNamedSkillMap[instance] = &interactiveElementNamedSkill{PacketId: InteractiveElementNamedSkillID}
	return interactiveElementNamedSkillMap[instance]
}

func (in *interactiveElementNamedSkill) Serialize(buff *bytes.Buffer) {
	in.interactiveElementSkill.Serialize(buff)
	utils.WriteVarInt32(buff, in.nameId)
}

func (in *interactiveElementNamedSkill) Deserialize(reader *bytes.Reader) {
	in.interactiveElementSkill = new(interactiveElementSkill)
	in.interactiveElementSkill.Deserialize(reader)
	in.nameId = utils.ReadVarInt32(reader)
}

func (in *interactiveElementNamedSkill) GetPacketId() uint32 {
	return in.PacketId
}

func (in *interactiveElementNamedSkill) String() string {
	return fmt.Sprintf("packetId: %d\n", in.PacketId)
}
