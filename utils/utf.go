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
	err := binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		panic(err)
	}

	bytesString := make([]byte, length)
	_, err = io.ReadFull(reader, bytesString)
	if err != nil {
		panic(err)
	}

	return string(bytesString)
}
