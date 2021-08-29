// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 09:23:06.4273962 +0200 CEST m=+0.020472201

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
	ancestors              []int32
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
	oEI                    []*objectEffectInteger
}

var mountClientMap = make(map[uint]*mountClient)

func GetMountClientNOA(instance uint) *mountClient {
	mountClient_, ok := mountClientMap[instance]

	if ok {
		return mountClient_
	}

	mountClientMap[instance] = &mountClient{PacketId: MountClientID}
	return mountClientMap[instance]
}

func (m *mountClient) Serialize(buff *bytes.Buffer) {
	var box uint32
	box = utils.SetFlag(box, 0, m.sex)
	box = utils.SetFlag(box, 1, m.isRideable)
	box = utils.SetFlag(box, 2, m.isWild)
	box = utils.SetFlag(box, 3, m.isFecondationReady)
	box = utils.SetFlag(box, 4, m.useHarnessColors)

	_ = binary.Write(buff, binary.BigEndian, byte(box))
	_ = binary.Write(buff, binary.BigEndian, m.id)
	utils.WriteVarInt32(buff, m.model)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(m.ancestors)))
	_ = binary.Write(buff, binary.BigEndian, m.ancestors)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(m.behaviors)))
	_ = binary.Write(buff, binary.BigEndian, m.behaviors)
	utils.WriteUTF(buff, m.name)
	_ = binary.Write(buff, binary.BigEndian, m.ownerId)
	utils.WriteVarLong(buff, m.experience)
	_ = binary.Write(buff, binary.BigEndian, m.experienceForNextLevel)
	_ = binary.Write(buff, binary.BigEndian, m.level)
	utils.WriteVarInt32(buff, m.maxPods)
	utils.WriteVarInt32(buff, m.stamina)
	utils.WriteVarInt32(buff, m.staminaMax)
	utils.WriteVarInt32(buff, m.maturity)
	utils.WriteVarInt32(buff, m.maturityForAdult)
	utils.WriteVarInt32(buff, m.energy)
	utils.WriteVarInt32(buff, m.energyMax)
	_ = binary.Write(buff, binary.BigEndian, m.serenity)
	_ = binary.Write(buff, binary.BigEndian, m.aggressivityMax)
	utils.WriteVarInt32(buff, m.serenityMax)
	utils.WriteVarInt32(buff, m.love)
	utils.WriteVarInt32(buff, m.loveMax)
	_ = binary.Write(buff, binary.BigEndian, m.fecondationTime)
	_ = binary.Write(buff, binary.BigEndian, m.boostLimiter)
	_ = binary.Write(buff, binary.BigEndian, m.boostMax)
	_ = binary.Write(buff, binary.BigEndian, m.reproductionCount)
	_ = binary.Write(buff, binary.BigEndian, m.reproductionCountMax)
	utils.WriteVarShort(buff, m.harnessGID)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(m.oEI)))
	for i := 0; i < len(m.oEI); i++ {
		m.oEI[i].Serialize(buff)
	}
}

func (m *mountClient) Deserialize(reader *bytes.Reader) {
	var box byte
	_ = binary.Read(reader, binary.BigEndian, &box)

	m.sex = utils.GetFlag(uint32(box), 0)
	m.isRideable = utils.GetFlag(uint32(box), 1)
	m.isWild = utils.GetFlag(uint32(box), 2)
	m.isFecondationReady = utils.GetFlag(uint32(box), 3)
	m.useHarnessColors = utils.GetFlag(uint32(box), 4)

	_ = binary.Read(reader, binary.BigEndian, &m.id)
	m.model = utils.ReadVarInt32(reader)
	var len0_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len0_)
	m.ancestors = make([]int32, len0_)
	_ = binary.Read(reader, binary.BigEndian, m.ancestors)
	var len2_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len2_)
	m.behaviors = make([]int32, len2_)
	_ = binary.Read(reader, binary.BigEndian, m.behaviors)
	m.name = utils.ReadUTF(reader)
	_ = binary.Read(reader, binary.BigEndian, &m.ownerId)
	m.experience = float64(utils.ReadVarUInt64(reader))
	_ = binary.Read(reader, binary.BigEndian, &m.experienceForNextLevel)
	_ = binary.Read(reader, binary.BigEndian, &m.level)
	m.maxPods = utils.ReadVarInt32(reader)
	m.stamina = utils.ReadVarInt32(reader)
	m.staminaMax = utils.ReadVarInt32(reader)
	m.maturity = utils.ReadVarInt32(reader)
	m.maturityForAdult = utils.ReadVarInt32(reader)
	m.energy = utils.ReadVarInt32(reader)
	m.energyMax = utils.ReadVarInt32(reader)
	_ = binary.Read(reader, binary.BigEndian, &m.serenity)
	_ = binary.Read(reader, binary.BigEndian, &m.aggressivityMax)
	m.serenityMax = utils.ReadVarInt32(reader)
	m.love = utils.ReadVarInt32(reader)
	m.loveMax = utils.ReadVarInt32(reader)
	_ = binary.Read(reader, binary.BigEndian, &m.fecondationTime)
	_ = binary.Read(reader, binary.BigEndian, &m.boostLimiter)
	_ = binary.Read(reader, binary.BigEndian, &m.boostMax)
	_ = binary.Read(reader, binary.BigEndian, &m.reproductionCount)
	_ = binary.Read(reader, binary.BigEndian, &m.reproductionCountMax)
	m.harnessGID = utils.ReadVarInt16(reader)
	var len_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len_)
	for i := 0; i < int(len_); i++ {
		m.oEI[i] = new(objectEffectInteger)
		m.oEI[i].Deserialize(reader)
	}
}

func (m *mountClient) GetPacketId() uint32 {
	return m.PacketId
}

func (m *mountClient) String() string {
	return fmt.Sprintf("packetId: %d\n", m.PacketId)
}
