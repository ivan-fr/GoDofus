// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 08:31:34.9423667 +0200 CEST m=+18.953475801

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type jobExperienceMultiUpdate struct {
	PacketId      uint32
	jobExperience []*jobExperience
}

var jobExperienceMultiUpdateMap = make(map[uint]*jobExperienceMultiUpdate)

func GetJobExperienceMultiUpdateNOA(instance uint) *jobExperienceMultiUpdate {
	jobExperienceMultiUpdate_, ok := jobExperienceMultiUpdateMap[instance]

	if ok {
		return jobExperienceMultiUpdate_
	}

	jobExperienceMultiUpdateMap[instance] = &jobExperienceMultiUpdate{PacketId: JobExperienceMultiUpdateID}
	return jobExperienceMultiUpdateMap[instance]
}

func (jo *jobExperienceMultiUpdate) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, uint16(len(jo.jobExperience)))
	for i := 0; i < len(jo.jobExperience); i++ {
		jo.jobExperience[i].Serialize(buff)
	}
}

func (jo *jobExperienceMultiUpdate) Deserialize(reader *bytes.Reader) {
	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	jo.jobExperience = nil
	for i := 0; i < int(len0_); i++ {
		aMessage0 := new(jobExperience)
		aMessage0.Deserialize(reader)
		jo.jobExperience = append(jo.jobExperience, aMessage0)
	}
}

func (jo *jobExperienceMultiUpdate) GetPacketId() uint32 {
	return jo.PacketId
}

func (jo *jobExperienceMultiUpdate) String() string {
	return fmt.Sprintf("packetId: %d\n", jo.PacketId)
}