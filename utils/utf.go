package utils

import (
	"bytes"
	"encoding/binary"
	"io"
)

func WriteUTF(buff *bytes.Buffer, value string) {
	bytesValue := []byte(value)
	_ = binary.Write(buff, binary.BigEndian, uint16(len(bytesValue)))
	_ = binary.Write(buff, binary.BigEndian, bytesValue)
}

func ReadUTF(reader *bytes.Reader) string {
	var length uint16
	_ = binary.Read(reader, binary.BigEndian, length)
	bytesString := make([]byte, length)
	_, _ = io.ReadFull(reader, bytesString)

	return string(bytesString)
}
