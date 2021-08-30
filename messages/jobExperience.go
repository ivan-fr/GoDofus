// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 08:29:37.3206265 +0200 CEST m=+92.950910701

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type jobExperience struct {
	PacketId            uint32
	jobId               byte
	jobLevel            byte
	jobXp               float64
	jobXpLevelFloor     float64
	jobXpNextLevelFloor float64
}

var jobExperienceMap = make(map[uint]*jobExperience)

func GetJobExperienceNOA(instance uint) *jobExperience {
	jobExperience_, ok := jobExperienceMap[instance]

	if ok {
		return jobExperience_
	}

	jobExperienceMap[instance] = &jobExperience{PacketId: JobExperienceID}
	return jobExperienceMap[instance]
}

func (jo *jobExperience) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, jo.jobId)
	_ = binary.Write(buff, binary.BigEndian, jo.jobLevel)
	utils.WriteVarInt64(buff, jo.jobXp)
	utils.WriteVarInt64(buff, jo.jobXpLevelFloor)
	utils.WriteVarInt64(buff, jo.jobXpNextLevelFloor)
}

func (jo *jobExperience) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &jo.jobId)
	_ = binary.Read(reader, binary.BigEndian, &jo.jobLevel)
	jo.jobXp = float64(utils.ReadVarInt64(reader))
	jo.jobXpLevelFloor = float64(utils.ReadVarInt64(reader))
	jo.jobXpNextLevelFloor = float64(utils.ReadVarInt64(reader))
}

func (jo *jobExperience) GetPacketId() uint32 {
	return jo.PacketId
}

func (jo *jobExperience) String() string {
	return fmt.Sprintf("packetId: %d\n", jo.PacketId)
}
