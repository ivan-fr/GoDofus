// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 10:05:46.9455779 +0200 CEST m=+0.021413501

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type spellList struct {
	PacketId              uint32
	spellPrevisualization bool
	sI                    []*spellItem
}

var spellListMap = make(map[uint]*spellList)

func GetSpellListNOA(instance uint) *spellList {
	spellList_, ok := spellListMap[instance]

	if ok {
		return spellList_
	}

	spellListMap[instance] = &spellList{PacketId: SpellListID}
	return spellListMap[instance]
}

func (s *spellList) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, s.spellPrevisualization)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(s.sI)))
	for i := 0; i < len(s.sI); i++ {
		s.sI[i].Serialize(buff)
	}
}

func (s *spellList) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &s.spellPrevisualization)
	var len_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len_)
	s.sI = nil
	for i := 0; i < int(len_); i++ {
		aSI := new(spellItem)
		aSI.Deserialize(reader)
		s.sI = append(s.sI, aSI)
	}
}

func (s *spellList) GetPacketId() uint32 {
	return s.PacketId
}

func (s *spellList) String() string {
	return fmt.Sprintf("packetId: %d\n", s.PacketId)
}
