// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-25 10:22:38.7155254 +0200 CEST m=+0.003141801

package messages

import (
	"GoDofus/settings"
	"GoDofus/utils"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
)

type selectedServerData struct {
	PacketId              uint32
	serverId              uint32
	Address               []byte
	portsLen              uint16
	Ports                 []int32
	canCreateNewCharacter bool
	ticket                []byte
	instance              uint
}

var selectedServerDataMap = make(map[uint]*selectedServerData)

func (s *selectedServerData) GetNOA(instance uint) Message {
	selectedServerData_, ok := selectedServerDataMap[instance]

	if ok {
		return selectedServerData_
	}

	selectedServerDataMap[instance] = &selectedServerData{PacketId: SelectedServerDataID, instance: instance}
	return selectedServerDataMap[instance]
}

func (s *selectedServerData) Serialize(buff *bytes.Buffer) {
	utils.WriteVarInt16(buff, int32(s.serverId))
	utils.WriteUTF(buff, []byte(settings.Settings.LocalAddress))
	_ = binary.Write(buff, binary.BigEndian, s.portsLen)

	for i := uint16(0); i < s.portsLen; i++ {
		utils.WriteVarInt16(buff, settings.Settings.LocalGamePort)
	}

	_ = binary.Write(buff, binary.BigEndian, s.canCreateNewCharacter)

	id := Types_[int(IdentificationID)].GetNOA(s.instance).(*Identification)
	theTicket := s.ticket
	aesKey := id.AesKEY_

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}

	if len(theTicket) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := aesKey[:aes.BlockSize]

	if len(theTicket)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(theTicket, theTicket)

	utils.WriteVarInt32(buff, int32(len(theTicket)))
	_ = binary.Write(buff, binary.BigEndian, theTicket)
}

func (s *selectedServerData) Deserialize(reader *bytes.Reader) {
	s.serverId = uint32(utils.ReadVarInt16(reader))
	s.Address = utils.ReadUTF(reader)
	_ = binary.Read(reader, binary.BigEndian, &s.portsLen)

	for i := uint16(0); i < s.portsLen; i++ {
		s.Ports = append(s.Ports, utils.ReadVarInt16(reader))
	}
	_ = binary.Read(reader, binary.BigEndian, &s.canCreateNewCharacter)

	lenght := utils.ReadVarInt32(reader)
	s.ticket = make([]byte, lenght)
	_ = binary.Read(reader, binary.BigEndian, s.ticket)
}

func (s *selectedServerData) GetPacketId() uint32 {
	return s.PacketId
}

func (s *selectedServerData) String() string {
	return fmt.Sprintf("PacketId: %d\nAddress: %s\nports: %v", s.PacketId, string(s.Address), s.Ports)
}
