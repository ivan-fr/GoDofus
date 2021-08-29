// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 22:29:26.1592243 +0200 CEST m=+62.097075801

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type itemsPreset struct {
	PacketId      uint32
	preset        *preset
	itemForPreset []*itemForPreset
	mountEquipped bool
	entityLook    *entityLook
}

var itemsPresetMap = make(map[uint]*itemsPreset)

func GetItemsPresetNOA(instance uint) *itemsPreset {
	itemsPreset_, ok := itemsPresetMap[instance]

	if ok {
		return itemsPreset_
	}

	itemsPresetMap[instance] = &itemsPreset{PacketId: ItemsPresetID}
	return itemsPresetMap[instance]
}

func (it *itemsPreset) Serialize(buff *bytes.Buffer) {
	it.preset.Serialize(buff)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(it.itemForPreset)))
	for i := 0; i < len(it.itemForPreset); i++ {
		it.itemForPreset[i].Serialize(buff)
	}
	_ = binary.Write(buff, binary.BigEndian, it.mountEquipped)
	it.entityLook.Serialize(buff)
}

func (it *itemsPreset) Deserialize(reader *bytes.Reader) {
	it.preset = new(preset)
	it.preset.Deserialize(reader)
	var len1_ uint16
	_ = binary.Read(reader, binary.BigEndian, &len1_)
	it.itemForPreset = nil
	for i := 0; i < int(len1_); i++ {
		aMessage1 := new(itemForPreset)
		aMessage1.Deserialize(reader)
		it.itemForPreset = append(it.itemForPreset, aMessage1)
	}
	_ = binary.Read(reader, binary.BigEndian, &it.mountEquipped)
	it.entityLook = new(entityLook)
	it.entityLook.Deserialize(reader)
}

func (it *itemsPreset) GetPacketId() uint32 {
	return it.PacketId
}

func (it *itemsPreset) String() string {
	return fmt.Sprintf("packetId: %d\n", it.PacketId)
}
