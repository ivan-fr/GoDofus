// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 14:47:26.2309037 +0200 CEST m=+53.105111001

package messages

import (
	"GoDofus/utils"
	"bytes"
	"fmt"
)

type achievementAchievedRewardable struct {
	PacketId            uint32
	achievementAchieved *achievementAchieved
	finishedlevel       int32
}

var achievementAchievedRewardableMap = make(map[uint]*achievementAchievedRewardable)

func (ac *achievementAchievedRewardable) GetNOA(instance uint) Message {
	achievementAchievedRewardable_, ok := achievementAchievedRewardableMap[instance]

	if ok {
		return achievementAchievedRewardable_
	}

	achievementAchievedRewardableMap[instance] = &achievementAchievedRewardable{PacketId: AchievementAchievedRewardableID}
	return achievementAchievedRewardableMap[instance]
}

func (ac *achievementAchievedRewardable) Serialize(buff *bytes.Buffer) {
	ac.achievementAchieved.Serialize(buff)
	utils.WriteVarInt16(buff, ac.finishedlevel)
}

func (ac *achievementAchievedRewardable) Deserialize(reader *bytes.Reader) {
	ac.achievementAchieved = new(achievementAchieved)
	ac.achievementAchieved.Deserialize(reader)
	ac.finishedlevel = utils.ReadVarInt16(reader)
}

func (ac *achievementAchievedRewardable) GetPacketId() uint32 {
	return ac.PacketId
}

func (ac *achievementAchievedRewardable) String() string {
	return fmt.Sprintf("packetId: %d\n", ac.PacketId)
}
