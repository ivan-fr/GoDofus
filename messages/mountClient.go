// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 21:19:05.6763423 +0200 CEST m=+343.489800401

package messages

import (
	"GoDofus/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

type mountClient struct {
	PacketId               uint32
	sex                    bool
	isRideable             bool
	isWild                 bool
	isFecondationReady     bool
	useHarnessColors       bool
	id                     float64
	model                  int32
	ancestor               []int32
	behaviors              []int32
	name                   []byte
	ownerId                int32
	experience             float64
	experienceForLevel     float64
	experienceForNextLevel float64
	level                  byte
	maxPods                int32
	stamina                int32
	staminaMax             int32
	maturity               int32
	maturityForAdult       int32
	energy                 int32
	energyMax              int32
	serenity               int32
	aggressivityMax        int32
	serenityMax            int32
	love                   int32
	loveMax                int32
	fecondationTime        int32
	boostLimiter           int32
	boostMax               float64
	reproductionCount      int32
	reproductionCountMax   int32
	harnessGID             int32
	objectEffectInteger    []*objectEffectInteger
}

var mountClientMap = make(map[uint]*mountClient)

func (mo *mountClient) GetNOA(instance uint) Message {
	mountClient_, ok := mountClientMap[instance]

	if ok {
		return mountClient_
	}

	mountClientMap[instance] = &mountClient{PacketId: MountClientID}
	return mountClientMap[instance]
}

func (mo *mountClient) Serialize(buff *bytes.Buffer) {
	var box0 uint32
	box0 = utils.SetFlag(box0, 0, mo.sex)
	box0 = utils.SetFlag(box0, 1, mo.isRideable)
	box0 = utils.SetFlag(box0, 2, mo.isWild)
	box0 = utils.SetFlag(box0, 3, mo.isFecondationReady)
	box0 = utils.SetFlag(box0, 4, mo.useHarnessColors)
	_ = binary.Write(buff, binary.BigEndian, byte(box0))
	_ = binary.Write(buff, binary.BigEndian, mo.id)
	utils.WriteVarInt32(buff, mo.model)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(mo.ancestor)))
	_ = binary.Write(buff, binary.BigEndian, mo.ancestor)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(mo.behaviors)))
	_ = binary.Write(buff, binary.BigEndian, mo.behaviors)
	utils.WriteUTF(buff, mo.name)
	_ = binary.Write(buff, binary.BigEndian, mo.ownerId)
	utils.WriteVarInt64(buff, mo.experience)
	utils.WriteVarInt64(buff, mo.experienceForLevel)
	_ = binary.Write(buff, binary.BigEndian, mo.experienceForNextLevel)
	_ = binary.Write(buff, binary.BigEndian, mo.level)
	utils.WriteVarInt32(buff, mo.maxPods)
	utils.WriteVarInt32(buff, mo.stamina)
	utils.WriteVarInt32(buff, mo.staminaMax)
	utils.WriteVarInt32(buff, mo.maturity)
	utils.WriteVarInt32(buff, mo.maturityForAdult)
	utils.WriteVarInt32(buff, mo.energy)
	utils.WriteVarInt32(buff, mo.energyMax)
	_ = binary.Write(buff, binary.BigEndian, mo.serenity)
	_ = binary.Write(buff, binary.BigEndian, mo.aggressivityMax)
	utils.WriteVarInt32(buff, mo.serenityMax)
	utils.WriteVarInt32(buff, mo.love)
	utils.WriteVarInt32(buff, mo.loveMax)
	_ = binary.Write(buff, binary.BigEndian, mo.fecondationTime)
	_ = binary.Write(buff, binary.BigEndian, mo.boostLimiter)
	_ = binary.Write(buff, binary.BigEndian, mo.boostMax)
	_ = binary.Write(buff, binary.BigEndian, mo.reproductionCount)
	utils.WriteVarInt32(buff, mo.reproductionCountMax)
	utils.WriteVarInt16(buff, mo.harnessGID)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(mo.objectEffectInteger)))
	for i := 0; i < len(mo.objectEffectInteger); i++ {
		mo.objectEffectInteger[i].Serialize(buff)
	}
}

func (mo *mountClient) Deserialize(reader *bytes.Reader) {
	var box0 byte
	_ = binary.Read(reader, binary.BigEndian, &box0)
	mo.sex = utils.GetFlag(uint32(box0), 0)
	mo.isRideable = utils.GetFlag(uint32(box0), 1)
	mo.isWild = utils.GetFlag(uint32(box0), 2)
	mo.isFecondationReady = utils.GetFlag(uint32(box0), 3)
	mo.useHarnessColors = utils.GetFlag(uint32(box0), 4)
	_ = binary.Read(reader, binary.BigEndian, &mo.id)
	mo.model = utils.ReadVarInt32(reader)
	var len3_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len3_)
	mo.ancestor = make([]int32, len3_)
	_ = binary.Read(reader, binary.BigEndian, mo.ancestor)
	var len4_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len4_)
	mo.behaviors = make([]int32, len4_)
	_ = binary.Read(reader, binary.BigEndian, mo.behaviors)
	mo.name = utils.ReadUTF(reader)
	_ = binary.Read(reader, binary.BigEndian, &mo.ownerId)
	mo.experience = float64(utils.ReadVarInt64(reader))
	mo.experienceForLevel = float64(utils.ReadVarInt64(reader))
	_ = binary.Read(reader, binary.BigEndian, &mo.experienceForNextLevel)
	_ = binary.Read(reader, binary.BigEndian, &mo.level)
	mo.maxPods = utils.ReadVarInt32(reader)
	mo.stamina = utils.ReadVarInt32(reader)
	mo.staminaMax = utils.ReadVarInt32(reader)
	mo.maturity = utils.ReadVarInt32(reader)
	mo.maturityForAdult = utils.ReadVarInt32(reader)
	mo.energy = utils.ReadVarInt32(reader)
	mo.energyMax = utils.ReadVarInt32(reader)
	_ = binary.Read(reader, binary.BigEndian, &mo.serenity)
	_ = binary.Read(reader, binary.BigEndian, &mo.aggressivityMax)
	mo.serenityMax = utils.ReadVarInt32(reader)
	mo.love = utils.ReadVarInt32(reader)
	mo.loveMax = utils.ReadVarInt32(reader)
	_ = binary.Read(reader, binary.BigEndian, &mo.fecondationTime)
	_ = binary.Read(reader, binary.BigEndian, &mo.boostLimiter)
	_ = binary.Read(reader, binary.BigEndian, &mo.boostMax)
	_ = binary.Read(reader, binary.BigEndian, &mo.reproductionCount)
	mo.reproductionCountMax = utils.ReadVarInt32(reader)
	mo.harnessGID = utils.ReadVarInt16(reader)
	var len29_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len29_)
	mo.objectEffectInteger = nil
	for i := 0; i < int(len29_); i++ {
		aMessage29 := new(objectEffectInteger)
		aMessage29.Deserialize(reader)
		mo.objectEffectInteger = append(mo.objectEffectInteger, aMessage29)
	}
}

func (mo *mountClient) GetPacketId() uint32 {
	return mo.PacketId
}

func (mo *mountClient) String() string {
	return fmt.Sprintf("packetId: %d\n", mo.PacketId)
}
