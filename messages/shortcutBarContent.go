// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-29 21:42:00.9066123 +0200 CEST m=+20.026969901

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type shortcutBarContent struct {
	PacketId    uint32
	barType     byte
	itemWrapper *itemWrapper
}

var shortcutBarContentMap = make(map[uint]*shortcutBarContent)

func GetShortcutBarContentNOA(instance uint) *shortcutBarContent {
	shortcutBarContent_, ok := shortcutBarContentMap[instance]

	if ok {
		return shortcutBarContent_
	}

	shortcutBarContentMap[instance] = &shortcutBarContent{PacketId: ShortcutBarContentID}
	return shortcutBarContentMap[instance]
}

func (sh *shortcutBarContent) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, sh.barType)
	sh.itemWrapper.Serialize(buff)
}

func (sh *shortcutBarContent) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &sh.barType)
	sh.itemWrapper = new(itemWrapper)
	sh.itemWrapper.Deserialize(reader)
}

func (sh *shortcutBarContent) GetPacketId() uint32 {
	return sh.PacketId
}

func (sh *shortcutBarContent) String() string {
	return fmt.Sprintf("packetId: %d\n", sh.PacketId)
}
