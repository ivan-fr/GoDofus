package utils

import (
	"bytes"
	"encoding/binary"
)

func WriteUTF(buff *bytes.Buffer, value string) {
	bytesValue := []byte(value)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(bytesValue)))
	_ = binary.Write(buff, binary.LittleEndian, bytesValue)
}
