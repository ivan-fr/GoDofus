// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-28 17:50:43.0200753 +0200 CEST m=+0.019985601

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type item struct {
	PacketId uint32
	content  []byte
	typeId   uint16
	myProcol Message
}

var protocolType = getProtocolType()

func (i *item) GetNOA(instance uint) Message {
	return nil
}

func getProtocolType() map[uint16]reflect.Type {

	var _typesTypes = make(map[uint16]reflect.Type)
	_typesTypes[ObjectEffectID] = reflect.TypeOf(objectEffect{})
	_typesTypes[ObjectEffectIntegerID] = reflect.TypeOf(objectEffectInteger{})
	_typesTypes[ObjectEffectCreatureID] = reflect.TypeOf(objectEffectCreature{})
	_typesTypes[ObjectEffectLadderID] = reflect.TypeOf(objectEffectLadder{})
	_typesTypes[ObjectEffectMinMaxID] = reflect.TypeOf(objectEffectMinMax{})
	_typesTypes[ObjectEffectDurationID] = reflect.TypeOf(objectEffectDuration{})
	_typesTypes[ObjectEffectStringID] = reflect.TypeOf(objectEffectString{})
	_typesTypes[ObjectEffectDiceID] = reflect.TypeOf(objectEffectDice{})
	_typesTypes[ObjectEffectDateID] = reflect.TypeOf(objectEffectDate{})
	_typesTypes[ObjectEffectMountID] = reflect.TypeOf(objectEffectMount{})
	_typesTypes[ItemsPresetID] = reflect.TypeOf(itemsPreset{})

	return _typesTypes
}

func (i *item) Serialize(buff *bytes.Buffer) {
	_ = binary.Write(buff, binary.BigEndian, i.typeId)
	i.myProcol.Serialize(buff)
}

func (i *item) Deserialize(reader *bytes.Reader) {
	_ = binary.Read(reader, binary.BigEndian, &i.typeId)
	pType, ok := protocolType[i.typeId]

	if !ok {
		panic(i.typeId)
	}

	newProtocol := reflect.New(pType).Interface().(Message)
	newProtocol.Deserialize(reader)
	i.myProcol = newProtocol
}

func (i *item) GetPacketId() uint32 {
	return i.PacketId
}

func (i *item) String() string {
	return fmt.Sprintf("packetId: %d\n", i.PacketId)
}
